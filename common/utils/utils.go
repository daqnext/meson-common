package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	sysnet "net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GetStringHash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//生成随机位数数字
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func FormatStrcut(s interface{}) string {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("%+v", s)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", s)
	}
	return out.String()
}

func RunCommand(cmdstring string, args ...string) (string, error) {
	cmd := exec.Command(cmdstring, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), err
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func GetDirSize(rootPath string) (uint64, error) {
	dirSize := uint64(0)

	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += uint64(file.Size())
		}

		return nil
	}

	err := filepath.Walk(rootPath, readSize)
	return dirSize, err
}

func GetMainMacAddress() (string, error) {
	ifas, err := sysnet.Interfaces()
	if err != nil {
		return "", err
	}

	ans := ""
	ansIndex := 1024

	for _, ifa := range ifas {
		// fmt.Printf("%+v %+v\n", ifa, uint(ifa.Flags))

		// Flags(19) means `up|broadcast|multicast`
		if ifa.Flags == sysnet.Flags(19) && ifa.Index < ansIndex {
			ans = ifa.HardwareAddr.String()
			ansIndex = ifa.Index
		}
	}
	return ans, nil
}

func HashBytes(input []byte) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func HashLocalFile(inputUrl string, bytesnum int64) (string, error) {

	f, err := os.Open(inputUrl)
	if err != nil {
		return "", err
	}

	btoready := make([]byte, bytesnum)
	n1, err := f.Read(btoready)
	if err != nil {
		return "", err
	}

	return HashBytes(btoready[:n1]), nil
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
