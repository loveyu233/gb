package snowflake

import (
	"log"
)

var worker *Worker

func init() {
	w, err := NewWorker(1)
	if err != nil {
		log.Fatal("雪花算法初始化失败")
		return
	}
	worker = w
}
func GetId() int64 {
	return worker.GetId()
}
