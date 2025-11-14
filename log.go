package gb

import "log"

type GBLog interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type GbDefaultlogger struct {
}

// init 函数用于处理init相关逻辑。
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// Infof 方法用于处理Infof相关逻辑。
func (l GbDefaultlogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Debugf 方法用于处理Debugf相关逻辑。
func (l GbDefaultlogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Errorf 方法用于处理Errorf相关逻辑。
func (l GbDefaultlogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
