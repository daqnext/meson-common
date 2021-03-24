package downloadtaskmgr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/daqnext/meson-common/common/logger"
	"github.com/daqnext/meson-common/common/utils"
)

type taskState string

const (
	Wait    taskState = "wait"
	Running taskState = "running"
	Finish  taskState = "finish"
)

type DownloadTask struct {
	Id           uint
	TargetUrl    string
	OriginTag    string
	BindNameHash string
	FileNameHash string
	TryTimes     int
	State        taskState // wait running finish
	Continent    string
	Country      string
	Area         string
}

var currentId uint
var idLock sync.Mutex
var leftTaskCount int
var runningTaskId uint

var execFunc func(task *DownloadTask) error
var onTaskFailed func(task *DownloadTask)

var fileLock sync.RWMutex
var fileHandle *os.File
var delTaskFileHandle *os.File
var recordFilePath string
var recordWriter *bufio.Writer

var jobChan = make(chan *DownloadTask, 1024)

var speedLimit = []int64{0, 10 * 1e3, 50 * 1e3, 200 * 1e3, 1000 * 1e3, 1000 * 1e6} //download channel speed line [0 KB/s, 10KB/s, 50KB/s, 200KB/s, 1000KB/s, 1000MB/s]
var countChan = make(chan bool, 5)

func InitTaskMgr(rootPath string) {
	if !utils.Exists(rootPath) {
		err := os.Mkdir(rootPath, 0700)
		if err != nil {
			logger.Fatal("tempfile dir create failed, please create dir " + rootPath + " by manual")
		}
	}

	recordFilePath = rootPath + "/unfinishtask.txt"
	fileHandle, err := os.OpenFile(recordFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("Open unfinishtask record file error", "err", err)
		return
	}

	defer func() {
		fileHandle.Close()
		fileHandle = nil
	}()

	//allow 5 tasks in same time
	for i := 0; i < 5; i++ {
		countChan <- true
	}

	//read unfinished task
	reader := bufio.NewReader(fileHandle)

	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Print(str)
		var task DownloadTask
		err = json.Unmarshal([]byte(str), &task)
		if err != nil {
			logger.Error("Unmarshal DownloadTask Error", "err", err)
			continue
		}
		currentId = task.Id
		leftTaskCount++
		jobChan <- &task
	}
}

func AddTask(targetUrl string, originTag string, continent string, country string, area string, bindNameHash string, fileNameHash string, tryTimes int) error {
	if fileHandle == nil {
		var err error
		fileHandle, err = os.OpenFile(recordFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			logger.Error("Open unfinishtask record file error", "err", err)
			return errors.New("open task record file error")
		}
		recordWriter = bufio.NewWriter(fileHandle)
	}
	idLock.Lock()
	currentId++
	idLock.Unlock()

	newTask := &DownloadTask{
		Id:           currentId,
		TargetUrl:    targetUrl,
		OriginTag:    originTag,
		BindNameHash: bindNameHash,
		FileNameHash: fileNameHash,
		TryTimes:     tryTimes,
		State:        Wait, // wait running finish
		Continent:    continent,
		Country:      country,
		Area:         area,
	}
	leftTaskCount++

	//add task to chan
	jobChan <- newTask

	//add task to the end of the file,the task can be continue when restart
	str, err := json.Marshal(*newTask)
	if err != nil {
		logger.Error("DownloadTask Marshal Error", "err", err)
		return errors.New("marshal downloadtask error")
	}
	recordWriter.WriteString(string(str) + "\n")

	fileLock.Lock()
	recordWriter.Flush()
	fileLock.Unlock()

	return nil
}

func SetExecTaskFunc(function func(task *DownloadTask) error) {
	execFunc = function
}

func SetOnTaskFailed(function func(task *DownloadTask)) {
	onTaskFailed = function
}

func Run() {
	go func() {
		for true {
			select {
			case task := <-jobChan:
				if execFunc == nil {
					logger.Error("execFun is nil, no func to exec task")
					return
				}
				go func() {
					<-countChan
					err := execFunc(task)
					countChan <- true
					if err != nil {
						//task failed
						//push to the end of the list
						if task.TryTimes < 3 {
							tryTimes := task.TryTimes + 1
							AddTask(task.TargetUrl, task.OriginTag, task.Continent, task.Country, task.Area, task.BindNameHash, task.FileNameHash, tryTimes)
						} else {
							if onTaskFailed != nil {
								onTaskFailed(task)
							}
						}

					}
					leftTaskCount--
					//delete task from record file
					RemoveFinishedTaskFromFile()
				}()

			}
		}
	}()
}

//delete first line of task list
func RemoveFinishedTaskFromFile() ([]byte, error) {
	if delTaskFileHandle == nil {
		var err error
		delTaskFileHandle, err = os.OpenFile(recordFilePath, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("open file error", err)
		}
	}
	f := delTaskFileHandle
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	line, err := buf.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	nw, err := io.Copy(f, buf)
	if err != nil {
		return nil, err
	}
	fileLock.Lock()
	defer fileLock.Unlock()

	err = f.Truncate(nw)
	if err != nil {
		return nil, err
	}
	err = f.Sync()
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return line, nil
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		if rwTimeout > 0 {
			err := conn.SetDeadline(time.Now().Add(rwTimeout))
			if err != nil {
				logger.Error("set download process rwTimeout error", "err", err)
				return nil, err
			}
		}
		return conn, nil
	}
}

func DownLoadFile(url string, distFilePath string) error {
	connectTimeout := 10 * time.Second
	//readWriteTimeout := 3600 * 3 * time.Second
	readWriteTimeout := time.Duration(0)
	//http client
	c := http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}

	//get
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("create request error", "err", err)
		return err
	}
	//下载文件
	response, err := c.Do(req)
	if err != nil {
		logger.Error("get file url "+url+" error", "err", err)
		return err
	}
	//creat folder and file
	distDir := path.Dir(distFilePath)
	err = os.MkdirAll(distDir, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(distFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if response.Body == nil {
		return errors.New("body is null")
	}
	defer response.Body.Close()

	//buff := make([]byte, 32*1024)
	//written := 0
	//writer := bufio.NewWriter(file)
	//reader := bufio.NewReaderSize(response.Body, 1024*32)
	//go func() {
	//	for {
	//		nr, er := reader.Read(buff)
	//		if nr > 0 {
	//			nw, ew := writer.Write(buff[0:nr])
	//			if nw > 0 {
	//				written += nw
	//			}
	//			if ew != nil {
	//				err = ew
	//				break
	//			}
	//			if nr != nw {
	//				err = io.ErrShortWrite
	//				break
	//			}
	//		}
	//		if er != nil {
	//			if er != io.EOF {
	//				err = er
	//			}
	//			break
	//		}
	//	}
	//	if err != nil {
	//		logger.Error("download file error","err",err)
	//	}
	//}()
	//
	//spaceTime := time.Second * 1
	//ticker := time.NewTicker(spaceTime)
	//lastWtn := 0
	//stop := false
	//
	//for {
	//	select {
	//	case <-ticker.C:
	//		speed := written - lastWtn
	//		fmt.Printf("[*] Speed %s / %s \n", speed, spaceTime.String())
	//		if written-lastWtn == 0 {
	//			ticker.Stop()
	//			stop = true
	//			break
	//		}
	//		lastWtn = written
	//	}
	//	if stop {
	//		break
	//	}
	//}

	//_, err = io.Copy(file, response.Body)

	//get a speed limit
	speedLimit := 9999999

	_, err = copyBuffer2(file, response.Body, nil, int64(speedLimit))

	if err != nil {
		os.Remove(distFilePath)
		return err
	}
	fileInfo, err := os.Stat(distFilePath)
	if err != nil {
		os.Remove(distFilePath)
		return err
	}
	size := fileInfo.Size()
	logger.Debug("donwload file,fileInfo", "size", size)
	//if size != length {
	//	os.Remove(distFilePath)
	//	return errors.New("download file size error")
	//}

	if size == 0 {
		os.Remove(distFilePath)
		return errors.New("download file size error")
	}

	return nil
}

func copyBuffer2(dst io.Writer, src io.Reader, buf []byte, speedLimit int64) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	//if wt, ok := src.(io.WriterTo); ok {
	//	return wt.WriteTo(dst)
	//}
	//// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	//if rt, ok := dst.(io.ReaderFrom); ok {
	//	return rt.ReadFrom(src)
	//}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	//written:=int64(0)
	stop := false
	go func() {
		for {
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if nw > 0 {
					written += int64(nw)
				}

				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
		}
		stop = true
	}()

	spaceTime := time.Millisecond * 200
	ticker := time.NewTicker(spaceTime)
	lastWtn := int64(0)
	lowSpeedCount := 0
	for {
		select {
		case <-ticker.C:
			speed := (written - lastWtn) * 5
			fmt.Printf("[*] Speed %d / %s \n", speed, spaceTime.String())
			if speed < speedLimit {
				lowSpeedCount++
			} else {
				lowSpeedCount = 0
			}
			if lowSpeedCount > 25 {
				ticker.Stop()
				logger.Error("download speed low", "speed", speed)
				return written, errors.New("slow download speed")
			}

			lastWtn = written
		}
		if stop {
			break
		}
	}

	return written, err
}

func copyBuffer(dst io.Writer, src io.Reader, buf []byte, speedLimit int64) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	//if wt, ok := src.(io.WriterTo); ok {
	//	return wt.WriteTo(dst)
	//}
	//// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	//if rt, ok := dst.(io.ReaderFrom); ok {
	//	return rt.ReadFrom(src)
	//}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}

	// startTime
	startTime := time.Now().UnixNano()
	usedTime := float64(0)
	speed := float64(0)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			usedTime = float64((time.Now().UnixNano() - startTime) / int64(time.Millisecond))
			if usedTime > 3000 {
				speed = float64(written) / float64(usedTime) // KB/s
				logger.Debug("download speed", "speed", speed)
				if speed < float64(speedLimit) {
					return written, errors.New("download speed slow")
				}
			}

			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
