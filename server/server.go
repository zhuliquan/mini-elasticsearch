package server

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	// server port
	addr string

	// cancel
	cancel context.CancelFunc

	// quit signal
	quit chan struct{}
}

func (s *Server) Run(ctx context.Context) error {
	router := gin.Default()
	s.Index(router)

	var listener net.Listener
	var err error
	if s.addr == "" {
		listener, err = net.Listen("tcp", ":0")
	} else {
		listener, err = net.Listen("tcp", s.addr)
	}
	if err != nil {
		return err
	}
	
	s.addr = listener.Addr().String()
	srv := &http.Server{
		Addr:    s.addr,
		Handler: router.Handler(),
	}
	go func() {
		_ = srv.Serve(listener)
		s.quit <- struct{}{}
		log.Println("exit from mini server")
	}()

	newCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	go func() {
		<-newCtx.Done()
		_ = srv.Shutdown(newCtx)
	}()

	return nil
}

func (s *Server) Close() {
	s.cancel()
}

func (s *Server) Quit() <-chan struct{} {
	return s.quit
}
