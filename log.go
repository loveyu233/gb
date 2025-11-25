package gb

import "log"

type GBLog interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type GbDefaultlogger struct {
}

// init 用来设置标准库日志的默认格式。
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// Infof 用来输出信息级日志。
func (l GbDefaultlogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Debugf 用来输出调试级日志。
func (l GbDefaultlogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Errorf 用来输出错误级日志。
func (l GbDefaultlogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
