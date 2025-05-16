package gb

import (
	"os"
	"os/signal"
	"syscall"
)

var _ Hook = (*SignalHook)(nil)

// Hook a graceful shutdown hook, default with signals of SIGINT and SIGTERM
type Hook interface {
	// WithSignals add more signals into hook
	WithSignals(signals ...syscall.Signal) Hook

	// Close register shutdown handles
	Close(funcs ...func())
}

type SignalHook struct {
	ctx chan os.Signal
}

// NewHook create a Hook instance
func NewHook() Hook {
	hook := &SignalHook{
		ctx: make(chan os.Signal, 1),
	}

	return hook.WithSignals(syscall.SIGINT, syscall.SIGTERM)
}

func (h *SignalHook) WithSignals(signals ...syscall.Signal) Hook {
	for _, s := range signals {
		signal.Notify(h.ctx, s)
	}

	return h
}

func (h *SignalHook) Close(funcs ...func()) {
	select {
	case <-h.ctx:
	}
	signal.Stop(h.ctx)

	for _, f := range funcs {
		f()
	}
}
