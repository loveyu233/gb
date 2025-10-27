package gb

import (
	"fmt"
	"mime/multipart"
	"sync"
)

func ReqKeywordAssembly(keyword string) string {
	return fmt.Sprintf("%%%s%%", keyword)
}

func ReqPageSize(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 20
	}
	return page, size
}

func ReqFileUploadGoroutine(files []*multipart.FileHeader, uploadFileFunc func(file *multipart.FileHeader) (string, error)) (fileURLS []string, errs []error) {
	var (
		group sync.WaitGroup
	)

	for _, ele := range files {
		group.Add(1)
		go func() {
			defer group.Done()
			url, err := uploadFileFunc(ele)
			if err != nil {
				errs = append(errs, err)
			} else {
				fileURLS = append(fileURLS, url)
			}
		}()
	}
	group.Wait()

	return fileURLS, errs
}
