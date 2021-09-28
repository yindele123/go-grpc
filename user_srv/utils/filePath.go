package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

func GetCurrentPath() string {

	file, _ := exec.LookPath(os.Args[0])

	//得到全路径，比如在windows下E:\\golang\\test\\a.exe
	path, _ := filepath.Abs(file)

	rst := filepath.Dir(path)
	return rst
}
