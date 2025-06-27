package gb

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetFileContentType 获取文件类型
func GetFileContentType(file []byte) string {
	return http.DetectContentType(file)
}

// ReadFileContent 读取文件内容
func ReadFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// GetCurrentLine 获取调用该方法的文件详细信息
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
