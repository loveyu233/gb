package gb

import (
	"fmt"
	"mime/multipart"
	"sync"
)

// ReqKeywordAssembly 用来把关键字包装成模糊查询格式。
func ReqKeywordAssembly(keyword string) string {
	return fmt.Sprintf("%%%s%%", keyword)
}

// ReqPageSize 用来校正页码、页大小并计算偏移量。
func ReqPageSize(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 20
	}
	return page, (page - 1) * size
}

// ReqFileUploadGoroutine 用来并发上传多文件并收集结果。
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
