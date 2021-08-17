package main

import (
	"context"
	"github.com/pengdafu/agent"
	"log"
	"net/http"
	"syscall"
	"time"
)

type server struct {
	svr http.Server
}

func (s *server) Start() error {
	return s.svr.ListenAndServe()
}

func (s *server) Stop() {
	ctx, fn := context.WithTimeout(context.Background(), time.Second*2)
	_ = s.svr.Shutdown(ctx)
	fn()
}

func (s *server) handlerHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func main() {
	svr := &server{}
	http.HandleFunc("/hello", svr.handlerHello)
	svr.svr = http.Server{Addr: "127.0.0.1:8080"}

	ag := agent.New()

	ag.RegisterCmd(svr)
	ag.HandlerSignals(syscall.SIGINT, syscall.SIGTERM)

	err := ag.Run()
	//err := svr.svr.ListenAndServe()
	log.Println("svr shutdown:", err)
}
