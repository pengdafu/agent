package agent

import (
	"fmt"
	"os"
	sysSignal "os/signal"
	"syscall"
)

type signal struct {
	signals []os.Signal
}

func (s signal) Start() error {
	ch := make(chan os.Signal, 1)
	sysSignal.Notify(ch, s.signals...)
	select {
	case sig := <-ch:
		return fmt.Errorf("recv signal %v", sig)
	}
}

func (s signal) Stop() {

}

func newSignalCmd(signals ...os.Signal) IAgent {
	if len(signals) == 0 {
		signals = append(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	}
	sig := &signal{signals: signals}
	return sig
}
