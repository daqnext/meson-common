package downloadtaskmgr

import (
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/daqnext/meson-common/common/logger"
)

type DownloadInfo struct {
	TargetUrl    string
	OriginTag    string
	BindName     string
	FileName     string
	Continent    string
	Country      string
	Area         string
	SavePath     string
	DownloadType string
	OriginRegion string
	TargetRegion string
}

type TaskStatus string

//const Task_Success TaskStatus = "success"
//const Task_Fail TaskStatus ="fail"
const Task_UnStart TaskStatus = "unstart"
const Task_Break TaskStatus = "break"
const Task_Downloading TaskStatus = "downloading"

type DownloadTask struct {
	DownloadInfo
	Id              uint64
	Status          TaskStatus
	FileSize        int64
	SpeedKBs        float64
	DownloadedSize  int64
	TryTimes        int
	StartTime       int64
	ZeroSpeedSec    int
	DownloadChannel *DownloadChannel
}

type TaskList struct {
	TaskInQueue []DownloadTask
}

var currentId uint64
var idLock sync.Mutex

const GlobalDownloadTaskChanSize = 1024 * 10

var globalDownloadTaskChan = make(chan *DownloadTask, GlobalDownloadTaskChanSize)

var onTaskSuccess func(task *DownloadTask)
var onTaskFailed func(task *DownloadTask)
var panicCatcher func()
var onDownloadStart func(task *DownloadTask)
var onDownloading func(task *DownloadTask, usedTimeSec int)

type ExecResult string

const Success ExecResult = "Success"
const Fail ExecResult = "Fail"
const Break ExecResult = "Break"

type DownloadChannel struct {
	SpeedLimitKBs           int64
	CountLimit              int
	RunningCountControlChan chan bool
	IdleChan                chan *DownloadTask
}

var DownloadingTaskMap sync.Map

var ChannelRunningSize = []int{10, 6, 3, 3, 2}
var channelArray = []*DownloadChannel{
	{SpeedLimitKBs: 30, CountLimit: ChannelRunningSize[0], RunningCountControlChan: make(chan bool, ChannelRunningSize[0]), IdleChan: make(chan *DownloadTask, 1024*5)},   //30KB/s
	{SpeedLimitKBs: 100, CountLimit: ChannelRunningSize[1], RunningCountControlChan: make(chan bool, ChannelRunningSize[1]), IdleChan: make(chan *DownloadTask, 1024*5)},  //100KB/s
	{SpeedLimitKBs: 500, CountLimit: ChannelRunningSize[2], RunningCountControlChan: make(chan bool, ChannelRunningSize[2]), IdleChan: make(chan *DownloadTask, 1024*5)},  //500KB/s
	{SpeedLimitKBs: 1500, CountLimit: ChannelRunningSize[3], RunningCountControlChan: make(chan bool, ChannelRunningSize[3]), IdleChan: make(chan *DownloadTask, 1024*3)}, //1500KB/s
	{SpeedLimitKBs: 2500, CountLimit: ChannelRunningSize[4], RunningCountControlChan: make(chan bool, ChannelRunningSize[4]), IdleChan: make(chan *DownloadTask, 1024*3)}, //2500KB/s
}

const NewRunningTaskCount = 7

var newRunningTaskControlChan = make(chan bool, NewRunningTaskCount)

func AddTaskToDownloadingMap(task *DownloadTask) {
	DownloadingTaskMap.Store(task.Id, task)
}

func DeleteDownloadingTask(taskid uint64) {
	DownloadingTaskMap.Delete(taskid)
}

type BySpeed []*DownloadTask

func (t BySpeed) Len() int           { return len(t) }
func (t BySpeed) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t BySpeed) Less(i, j int) bool { return t[i].SpeedKBs < t[j].SpeedKBs }
func LoopScanRunningTask() {
	newWaitingTaskCount := len(globalDownloadTaskChan)
	//logger.Debug("Download waiting len","len",newWaitingTaskCount)
	if newWaitingTaskCount <= 0 {
		//logger.Debug("have no new task waiting")
		return
	}

	killCount := 3
	if newWaitingTaskCount < killCount {
		killCount = newWaitingTaskCount
	}

	nowTime := time.Now().Unix()
	taskReadyToKill := []*DownloadTask{}
	DownloadingTaskMap.Range(func(key, value interface{}) bool {
		task, ok := value.(*DownloadTask)
		if ok {
			//loop find task to kill
			if nowTime-task.StartTime < 5 {
				return true
			}

			for _, v := range channelArray {
				if task.SpeedKBs < float64(v.SpeedLimitKBs) {
					//if task.FileSize>0  {
					//	finishPercent:=task.DownloadedSize*100/task.FileSize
					//	if finishPercent>70 {
					//		break
					//	}
					//}
					task.DownloadChannel = v
					taskReadyToKill = append(taskReadyToKill, task)
					break
				}
			}
		}
		return true
	})

	if len(taskReadyToKill) == 0 {
		return
	}

	sort.Sort(BySpeed(taskReadyToKill))
	count := 0
	for _, v := range taskReadyToKill {
		v.Status = Task_Break
		//logger.Debug("Break Task","id",v.Id)
		count++
		if count >= killCount {
			return
		}
	}
}

func InitTaskMgr(rootPath string) {
	LevelDBInit()

	for _, v := range channelArray {
		for i := 0; i < v.CountLimit; i++ {
			v.RunningCountControlChan <- true
		}
	}

	for i := 0; i < NewRunningTaskCount; i++ {
		newRunningTaskControlChan <- true
	}

	//read unfinished task and restart
	unFinishedTask := LoopTasksInLDB()
	if unFinishedTask == nil {
		return
	}

	for _, v := range unFinishedTask {
		info := &DownloadInfo{}
		info.TargetUrl = v.TargetUrl
		info.OriginTag = v.OriginTag
		info.BindName = v.BindName
		info.FileName = v.FileName
		info.Continent = v.Continent
		info.Country = v.Country
		info.Area = v.Area
		info.SavePath = v.SavePath

		err := AddGlobalDownloadTask(info)
		if err != nil {
			logger.Error("Add AddGlobalDownloadTask error")
		}
	}
}

func AddGlobalDownloadTask(info *DownloadInfo) error {

	idLock.Lock()
	if currentId >= math.MaxUint64 {
		currentId = 0
	}
	currentId++
	idLock.Unlock()

	newTask := &DownloadTask{}
	newTask.Id = currentId
	newTask.TargetUrl = info.TargetUrl
	newTask.OriginTag = info.OriginTag
	newTask.BindName = info.BindName
	newTask.FileName = info.FileName
	newTask.Continent = info.Continent
	newTask.Country = info.Country
	newTask.Area = info.Area
	newTask.DownloadType = info.DownloadType
	newTask.OriginRegion = info.OriginRegion
	newTask.TargetRegion = info.TargetRegion
	newTask.SavePath = info.SavePath
	newTask.Status = Task_UnStart
	newTask.TryTimes = 0

	go func() {
		//save to LevelDB
		SetTaskToLDB(newTask)
		//to task channel
		globalDownloadTaskChan <- newTask
	}()

	return nil
}

func SetPanicCatcher(function func()) {
	panicCatcher = function
}

func SetOnTaskSuccess(function func(task *DownloadTask)) {
	onTaskSuccess = function
}

func SetOnTaskFailed(function func(task *DownloadTask)) {
	onTaskFailed = function
}

func SetOnDownloading(function func(task *DownloadTask, usedTimeSec int)) {
	onDownloading = function
}

func SetOnDownloadStart(function func(task *DownloadTask)) {
	onDownloadStart = function
}

func GetDownloadTaskList() []*DownloadTask {
	taskInLDB := LoopTasksInLDB()
	if taskInLDB == nil {
		return nil
	}

	list := []*DownloadTask{}
	for _, v := range taskInLDB {
		list = append(list, v)
	}
	return list
}

func TaskSuccess(task *DownloadTask) {
	logger.Debug("Task Success", "id", task.Id)
	//从map中删除任务
	DelTaskFromLDB(task.Id)
	DeleteDownloadingTask(task.Id)
	if onTaskSuccess == nil {
		logger.Error("not define onTaskSuccess")
		return
	}
	onTaskSuccess(task)
}

func TaskFail(task *DownloadTask) {
	logger.Debug("Task Fail", "id", task.Id)
	//从map中删除任务
	DelTaskFromLDB(task.Id)
	DeleteDownloadingTask(task.Id)
	if onTaskFailed == nil {
		logger.Error("not define onTaskFailed")
		return
	}
	onTaskFailed(task)
}

func TaskBreak(task *DownloadTask) {
	logger.Debug("Task Break", "id", task.Id)
	//delete from runningMap
	DeleteDownloadingTask(task.Id)
	task.Status = Task_UnStart
	//add to queue
	channel := task.DownloadChannel
	if channel == nil {
		logger.Error("Break Task not set channel,back to global list", "taskid", task.Id)
		globalDownloadTaskChan <- task
		return
	}
	channel.IdleChan <- task
	logger.Debug("add break task to idleChan", "speedLimit", channel.SpeedLimitKBs, "chanLen", len(channel.IdleChan), "taskid", task.Id)
}

func TaskRetry(task *DownloadTask) {
	logger.Debug("Task Retry", "id", task.Id)
	DeleteDownloadingTask(task.Id)
	task.TryTimes++
	task.Status = Task_UnStart
	globalDownloadTaskChan <- task
}

func StartTask(task *DownloadTask) {
	if panicCatcher != nil {
		defer panicCatcher()
	}

	result := ExecDownloadTask(task)
	switch result {
	case Success:
		//logger.Debug("download task success", "id", task.Id)
		TaskSuccess(task)
	case Fail:
		//logger.Debug("download task fail", "id", task.Id)
		if task.TryTimes > 3 {
			TaskFail(task)
		} else {
			//继续放入任务队列
			TaskRetry(task)
		}
	case Break:
		//logger.Debug("download task idle", "id", task.Id)
		TaskBreak(task)
	}
}

func (dc *DownloadChannel) ChannelDownload() {
	go func() {
		for true {
			//拿到自己队列的token
			<-dc.RunningCountControlChan
			select {
			case task := <-dc.IdleChan:
				go func() {
					defer func() {
						dc.RunningCountControlChan <- true
					}()
					logger.Debug("get a task from idle list", "channel speed", dc.SpeedLimitKBs, "id", task.Id, "chanlen", len(dc.IdleChan))
					//执行任务
					StartTask(task)
				}()
			}
		}
	}()
}

func Run() {
	RunNewTask()
	RunChannelDownload()

	//scanloop
	go func() {
		if panicCatcher != nil {
			defer panicCatcher()
		}
		for true {
			time.Sleep(5 * time.Second)
			LoopScanRunningTask()
		}
	}()
}
func RunChannelDownload() {
	for _, v := range channelArray {
		v.ChannelDownload()
	}
}
func RunNewTask() {
	go func() {
		for true {
			<-newRunningTaskControlChan
			select {
			case task := <-globalDownloadTaskChan:
				//开始一个新下载任务
				go func() {
					//任务结束,放回token
					defer func() {
						newRunningTaskControlChan <- true
					}()
					//执行任务
					//logger.Debug("start a new task", "id", task.Id)
					task.Status = Task_Downloading
					AddTaskToDownloadingMap(task)
					StartTask(task)
				}()
			}
		}
	}()
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

func ExecDownloadTask(task *DownloadTask) ExecResult {
	connectTimeout := 10 * time.Second
	readWriteTimeout := 3600 * 12 * time.Second
	//readWriteTimeout := time.Duration(0)

	url := task.TargetUrl
	distFilePath := task.SavePath

	cHead := http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
	//get
	reqHead, err := http.NewRequest(http.MethodHead, url, nil)
	if err == nil {
		responseHead, err := cHead.Do(reqHead)
		if err == nil {
			if responseHead.StatusCode == 200 && responseHead.ContentLength > 0 {
				task.FileSize = responseHead.ContentLength
			}
		}
	}

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
		return Fail
	}
	//download
	response, err := c.Do(req)
	if err != nil {
		logger.Error("get file url "+url+" error", "err", err)
		return Fail
	}
	if response.StatusCode != 200 {
		logger.Error("get file url "+url+" error", "err", err, "statusCode", response.StatusCode)
		return Fail
	}
	//creat folder and file
	distDir := path.Dir(distFilePath)
	err = os.MkdirAll(distDir, os.ModePerm)
	if err != nil {
		return Fail
	}
	file, err := os.Create(distFilePath)
	if err != nil {
		return Fail
	}
	defer file.Close()
	if response.Body == nil {
		logger.Error("Download responseBody is null")
		return Fail
	}
	defer response.Body.Close()

	task.StartTime = time.Now().Unix()
	if onDownloadStart != nil {
		go onDownloadStart(task)
	}

	_, err = copyBuffer(file, response.Body, nil, task)

	if err != nil {
		os.Remove(distFilePath)
		if err.Error() == string(Break) {
			//logger.Debug("task break","id",task.Id)
			return Break
		}
		return Fail
	}
	fileInfo, err := os.Stat(distFilePath)
	if err != nil {
		logger.Error("Get file Stat error", "err", err)
		os.Remove(distFilePath)
		return Fail
	}
	size := fileInfo.Size()
	logger.Debug("donwload file,fileInfo", "size", size)

	if size == 0 {
		os.Remove(distFilePath)
		logger.Error("download file size error")
		return Fail
	}

	return Success
}

func copyBuffer(dst io.Writer, src io.Reader, buf []byte, task *DownloadTask) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
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
	stop := false

	srcWithCloser, ok := src.(io.ReadCloser)
	if ok == false {
		err = errors.New("to io.ReadCloser error")
		return written, err
	}
	go func() {
		for {
			time.Sleep(500 * time.Millisecond) //for test
			nr, er := srcWithCloser.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if nw > 0 {
					written += int64(nw)
				}

				if ew != nil {
					err = ew
					//fmt.Println(ew.Error())
					if task.Status == Task_Break &&
						(strings.Contains(err.Error(), "http: read on closed response body") ||
							strings.Contains(err.Error(), "use of closed network connection")) {
						err = errors.New(string(Break))
					}
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
					//errStr:=err.Error()
					//fmt.Println(errStr)
					if task.Status == Task_Break &&
						(strings.Contains(err.Error(), "http: read on closed response body") ||
							strings.Contains(err.Error(), "use of closed network connection")) {
						err = errors.New(string(Break))
					}
				}
				break
			}
		}
		stop = true
	}()

	//monitor download speed
	spaceTime := time.Millisecond * 1000
	ticker := time.NewTicker(spaceTime)
	//lastWtn := int64(0)
	count := 0
	for {
		count++
		if stop {
			break
		}
		select {
		case <-ticker.C:
			if task.Status == Task_Break {
				srcWithCloser.Close()
			}

			task.DownloadedSize = written
			useTime := count * 1000
			speed := float64(written) / float64(useTime)
			task.SpeedKBs = speed
			//reportDownloadState
			if onDownloading != nil {
				go onDownloading(task, useTime)
			}

		}
	}
	return written, err
}
