package gb

import (
	"fmt"
	"mime/multipart"
	"sync"
)

// ReqKeywordAssembly 函数用于处理ReqKeywordAssembly相关逻辑。
func ReqKeywordAssembly(keyword string) string {
	return fmt.Sprintf("%%%s%%", keyword)
}

// ReqPageSize 函数用于处理ReqPageSize相关逻辑。
func ReqPageSize(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 20
	}
	return page, (page - 1) * size
}

// ReqFileUploadGoroutine 函数用于处理ReqFileUploadGoroutine相关逻辑。
func ReqFileUploadGoroutine(files []*multipart.FileHeader, uploadFileFunc func(file *multipart.FileHeader) (string, error)) (fileURLS []string, errs []error) {
	var (
		group sync.WaitGroup
		mu    sync.Mutex
	)

	for _, fileHeader := range files {
		fh := fileHeader
		group.Add(1)
		go func(header *multipart.FileHeader) {
			defer group.Done()
			url, err := uploadFileFunc(header)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err)
				return
			}
			fileURLS = append(fileURLS, url)
		}(fh)
	}
	group.Wait()

	return fileURLS, errs
}
