package logger

import (
	"errors"
	"fmt"
	"github.com/daqnext/meson-common/common/utils"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type LogFileWriter struct {
	RootDir string
	MaxSize int64
	file    *os.File
	//write count
	size int64
	//date
	lastDate        string
	lastCount       int
	OnLogFileChange func(fileName string)
}

// SliceToFields Convert the slice to logrus.Fields
func SliceToFields(params []interface{}) logrus.Fields {
	if len(params)%2 != 0 {
		logrus.Error("Log parameter length is wrong!")
		return nil
	}
	fields := logrus.Fields{}
	for i := 0; i < len(params); i += 2 {
		val := params[i]
		key, found := val.(string)
		if found {
			fields[key] = params[i+1]
		} else {
			return nil
		}
	}
	return fields
}

func (p *LogFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.RootDir == "" {
		p.RootDir = "./"
	}
	if p.MaxSize == 0 {
		p.MaxSize = 1024 * 512 //512k
	}
	if p.file == nil {
		p.lastDate = time.Now().Format("2006-01-02")

		//check today file is exist or not
		count := 0
		rd, err := ioutil.ReadDir(p.RootDir + "log")
		if err != nil {
			//serverlogger.Println("read dir err",err)
			err := os.Mkdir(p.RootDir+"log", 0777)
			if err != nil {
				fmt.Println(err)
				return 0, errors.New("Mkdir " + p.RootDir + "log error")
			}
		} else {
			for _, fi := range rd {
				if fi.IsDir() {
				} else {
					filename := fi.Name()
					if strings.Contains(filename, p.lastDate) {
						count++
					}
				}
			}
		}

		if count == 0 {
			p.lastCount = 1
		} else {
			p.lastCount = count
		}

		//open log file
		p.file, err = os.OpenFile(p.RootDir+"log"+"/"+p.lastDate+"-"+fmt.Sprintf("%04d", p.lastCount)+".log",
			os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0777)
		if err != nil {
			fmt.Println(err)
			return 0, errors.New("OpenFile " + p.RootDir + "log" + "/" + p.lastDate + "-" + fmt.Sprintf("%04d", p.lastCount) + ".log error")
		}
		info, err := p.file.Stat()
		if err != nil {
			fmt.Println(err)
			return 0, errors.New("fileStat " + p.RootDir + "log" + "/" + p.lastDate + "-" + fmt.Sprintf("%04d", p.lastCount) + ".log error")
		}
		p.size = info.Size()
	}
	n, e := p.file.Write(data)
	p.size += int64(n)
	//max size 512K byte
	if p.size > p.MaxSize {
		oldFileName := p.file.Name()
		p.file.Close()
		if p.OnLogFileChange != nil {
			p.OnLogFileChange(oldFileName)
		}
		//fmt.Println("log file full")

		p.file, err = os.OpenFile(p.RootDir+"log"+"/"+p.lastDate+"-"+fmt.Sprintf("%04d", p.lastCount+1)+".log",
			os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0777)
		if err != nil {
			fmt.Println(err)
			return 0, errors.New("OpenFile " + p.RootDir + "log" + "/" + p.lastDate + "-" + fmt.Sprintf("%04d", p.lastCount+1) + ".log error")
		}
		p.size = 0
		p.lastCount++
	}
	if time.Now().Format("2006-01-02") != p.lastDate {
		oldFileName := p.file.Name()
		p.file.Close()
		if p.OnLogFileChange != nil {
			p.OnLogFileChange(oldFileName)
		}
		//fmt.Println("log file date change")
		p.lastDate = time.Now().Format("2006-01-02")
		p.lastCount = 1
		p.file, err = os.OpenFile(p.RootDir+"log"+"/"+p.lastDate+"-"+fmt.Sprintf("%04d", p.lastCount)+".log",
			os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0777)
		if err != nil {
			fmt.Println(err)
			return 0, errors.New("OpenFile " + p.RootDir + "log" + "/" + p.lastDate + "-" + fmt.Sprintf("%04d", p.lastCount) + ".log error")
		}
		p.size = 0
	}
	return n, e
}

func DeleteLog(path string, passTimeSec int64) {
	nowTime := time.Now().Unix()
	//default log
	deleteFileNames := []string{}
	if !utils.Exists(path) {
		Debug("DeleteLog folder not exist", "path", path)
		return
	}
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		Error("read dir fail", "err", err, "dir", path)
		return
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			name := fi.Name()
			length := 10
			if len(name) < 10 {
				length = len(name)
			}
			dateStr := name[:length]
			timeStamp, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				Error("DeleteTimeoutLog time.Parse error", "err", err)
				continue
			}
			timeStampFile := timeStamp.Unix()
			if nowTime-timeStampFile > passTimeSec {
				deleteFileNames = append(deleteFileNames, name)
			}
		}
	}

	for _, v := range deleteFileNames {
		fileName := path + "/" + v
		os.Remove(fileName)
	}
}
