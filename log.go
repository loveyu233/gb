package gb

import "log"

type GBLog interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type GbDefaultlogger struct {
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (l GbDefaultlogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l GbDefaultlogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l GbDefaultlogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
