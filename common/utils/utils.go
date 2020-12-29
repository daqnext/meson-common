package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	sysnet "net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func GetStringHash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//rand num string of given length
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

// is folder or file exist
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

// is folder or not
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// is file or not
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

func findEmptyFolder(dirname string) (emptys []string, err error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return []string{dirname}, nil
	}

	for _, file := range files {
		if file.IsDir() {
			edirs, err := findEmptyFolder(path.Join(dirname, file.Name()))
			if err != nil {
				return nil, err
			}
			if edirs != nil {
				emptys = append(emptys, edirs...)
			}
		}
	}
	return emptys, nil
}

func DeleteEmptyFolders(path string) {
	emptys, err := findEmptyFolder(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, dir := range emptys {
		if err := os.Remove(dir); err != nil {
			fmt.Println("delete folder error:", err.Error())
		} else {

		}
	}
}

func FileAddMark(originFileName string, mark string) string {
	extension := filepath.Ext(originFileName)
	fileNameAddMark := ""
	if extension == "" {
		fileNameAddMark = originFileName + mark
	} else {
		fileNameAddMark = originFileName[:len(originFileName)-len(extension)]
		fileNameAddMark = fileNameAddMark + mark + extension
	}
	return fileNameAddMark
}

//compare version (x.x.x)
//
//a>b return 1
//
//a==b return 0
//
//a<b  return -1
func VersionCompare(a string, b string) int {
	aVersion := strings.Split(a, ".")
	bVersion := strings.Split(b, ".")
	for i, _ := range aVersion {
		if aVersion[i] > bVersion[i] {
			return 1
		} else if aVersion[i] < bVersion[i] {
			return -1
		}
	}
	return 0
}
