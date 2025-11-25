package gb

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetFileContentType 用来根据文件前缀识别内容类型。
func GetFileContentType(file []byte) string {
	return http.DetectContentType(file)
}

// GetFileNameType 用来取得文件名的扩展名。
func GetFileNameType(fileName string) string {
	split := strings.Split(fileName, ".")
	if len(split) < 2 {
		return ""
	}
	return split[len(split)-1]
}

// ReadFileContent 用来读取指定路径的文件内容。
func ReadFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// GetCurrentLine 用来返回当前调用点的文件、函数与行号。
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
