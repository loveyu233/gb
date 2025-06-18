package gb

import "os"

func ReadFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
