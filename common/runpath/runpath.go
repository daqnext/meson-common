package runpath

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var RunPath string

func init() {
	file, _ := exec.LookPath(os.Args[0])
	runPath, _ := filepath.Abs(file)
	index := strings.LastIndex(runPath, string(os.PathSeparator))
	RunPath = runPath[:index]
}
