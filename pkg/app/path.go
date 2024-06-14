package app

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	workDir = ""
	rootDir = ""
	tempDir = os.TempDir()
)

func rootPath() string {
	dir := ""
	exe, err := os.Executable()
	if err == nil {
		path := filepath.Base(exe)
		if !strings.HasPrefix(exe, tempDir) && !strings.HasPrefix(path, "___") {
			dir = filepath.Dir(exe)
		} else {
			_, filename, _, ok := runtime.Caller(0)
			if ok {
				// 需要根据当前文件所处目录，修改相对位置
				dir = filepath.Join(filepath.Dir(filename), "../../")
			}
		}
	}
	dir = filepath.Join(dir, workDir)
	return dir
}

func Wd() string {
	return rootDir
}

func Init(dir ...string) {
	workDir = "temp"
	if len(dir) > 0 {
		workDir = dir[0]
	}
	err := os.Chdir(rootPath())
	if err != nil {
		panic(err)
	}
}
