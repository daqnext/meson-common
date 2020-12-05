package downloadtaskmgr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daqnext/meson-common/common/logger"
	"github.com/daqnext/meson-common/common/utils"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
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
				err := execFunc(task)
				if err != nil {
					//task failed
					//push to the end of the list
					if task.TryTimes < 5 {
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

func DownLoadFile(url string, distFilePath string) error {
	//http client
	client := new(http.Client)
	//head, err := client.Head(url)
	//if err != nil {
	//	return err
	//}
	//length := head.ContentLength
	//logger.Debug("donwload file,head contentLength", "length", length)

	//get
	resp, err := client.Get(url)
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
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
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
