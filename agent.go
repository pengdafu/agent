package agent

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"os"
)

var (
	ErrNotCmd = errors.New("no register cmd")
)

type Agent struct {
	cmds    []IAgent
	signals []os.Signal
	g       *errgroup.Group
	ctx     context.Context
}

func New() *Agent {
	agent := &Agent{}
	agent.cmds = make([]IAgent, 0)
	agent.g, agent.ctx = errgroup.WithContext(context.Background())
	return &Agent{}
}

func (agent *Agent) Run() error {
	if len(agent.cmds) == 0 {
		return ErrNotCmd
	}

	agent.cmds = append(agent.cmds, newSignalCmd(agent.signals...))

	_ = agent.start()
	err := agent.g.Wait()
	agent.stop()

	return err
}

func (agent *Agent) RegisterCmd(cmd IAgent) {
	agent.cmds = append(agent.cmds, cmd)
}

func (agent *Agent) HandlerSignals(signals ...os.Signal) {
	agent.signals = signals
}

func (agent *Agent) start() error {
	for _, cmd := range agent.cmds {
		agent.g.Go(func() error {
			select {
			case err := <-agent.warpCmdStart(cmd.Start):
				return err
			case <-agent.ctx.Done():
				return agent.ctx.Err()
			}
		})
	}
	return nil
}

func (agent *Agent) warpCmdStart(fn func() error) <-chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- fn()
	}()
	return ch
}

func (agent *Agent) stop() {
	for _, cmd := range agent.cmds {
		cmd.Stop()
	}
}
