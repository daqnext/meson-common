package logger

import (
	"errors"
	"fmt"
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
	lastDate  string
	lastCount int
}

// Convert the slice to logrus.Fields
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
		p.file, _ = os.OpenFile(p.RootDir+"log"+"/"+p.lastDate+"-"+fmt.Sprintf("%04d", p.lastCount)+".log",
			os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0777)
		info, _ := p.file.Stat()
		p.size = info.Size()
	}
	n, e := p.file.Write(data)
	p.size += int64(n)
	//max size 512K byte
	if p.size > p.MaxSize {
		p.file.Close()
		fmt.Println("log file full")
		p.file, _ = os.OpenFile(p.RootDir+"log"+"/"+p.lastDate+"-"+fmt.Sprintf("%04d", p.lastCount)+".log",
			os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0777)
		p.size = 0
		p.lastCount++
	}
	if time.Now().Format("2006-01-02") != p.lastDate {
		p.file.Close()
		fmt.Println("log file date change")
		p.lastDate = time.Now().Format("2006-01-02")
		p.lastCount = 1
		p.file, _ = os.OpenFile(p.RootDir+"log"+"/"+p.lastDate+"-"+fmt.Sprintf("%04d", p.lastCount)+".log",
			os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0777)
		p.size = 0
	}
	return n, e
}
