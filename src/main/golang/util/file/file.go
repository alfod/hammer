package file

import (
	"log"
	"path/filepath"
	"runtime"
)

func GetCurrentFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	log.Panicln(filename)
	dir1, _ := filepath.Split(filename)
	return dir1
}
