package common

import (
	"path"
	"runtime"
)

var RootPath string //Project root path

func init() {
	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)
	RootPath = path.Join(dir, "/..")
}

const FileTransferRunPort = "9092"
const ValidatorRunPort = "9093"
const SpeedTesterRunPort = "9094"

const RedirectMark = "-redirecter456gt"
