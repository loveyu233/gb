package gb

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetFileContentType 函数用于处理GetFileContentType相关逻辑。
func GetFileContentType(file []byte) string {
	return http.DetectContentType(file)
}

// GetFileNameType 函数用于处理GetFileNameType相关逻辑。
func GetFileNameType(fileName string) string {
	split := strings.Split(fileName, ".")
	if len(split) < 2 {
		return ""
	}
	return split[len(split)-1]
}

// ReadFileContent 函数用于处理ReadFileContent相关逻辑。
func ReadFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// GetCurrentLine 函数用于处理GetCurrentLine相关逻辑。
func GetCurrentLine() (folderName string, fileName string, funcName string, lineNumber int) {
	pc, file, line, _ := runtime.Caller(1)
	lineNumber = line
	funcName = runtime.FuncForPC(pc).Name()
	if idx := strings.LastIndex(funcName, "."); idx != -1 {
		funcName = funcName[idx+1:]
	}

	fileName = filepath.Base(file)
	folderName = filepath.Base(filepath.Dir(file))

	return
}
