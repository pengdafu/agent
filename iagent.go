package agent

type IAgent interface {
	Start() error
	Stop()
}
