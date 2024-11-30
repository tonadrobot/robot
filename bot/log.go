package bot

import (
	"fmt"
	"log"
	"path"
	"runtime"
)

// Logs error
func loge(err error) {
	errMsg := getCallerInfo() + err.Error()
	notify(errMsg, Admin)
	log.Println(errMsg)
}

// Logs string
func logs(errMsg string) {
	errMsg = getCallerInfo() + errMsg
	notify(errMsg, Admin)
	log.Println(errMsg)
}

// Returns file and line numbef of the caller of some function
func getCallerInfo() (info string) {

	// pc, file, lineNo, ok := runtime.Caller(2)
	_, file, lineNo, ok := runtime.Caller(2)
	if !ok {
		info = "runtime.Caller() failed"
		return
	}
	// funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // The Base function returns the last element of the path
	return fmt.Sprintf("%s:%d: ", fileName, lineNo)
}
