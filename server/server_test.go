package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("test exit with context.Done", func(t *testing.T) {
		st := time.Now()

		// set run time
		runTime := time.Second * 3

		// run srv
		srv := Server{quit: make(chan struct{}, 1)}
		ctx, cancel := context.WithTimeout(context.Background(), runTime)
		defer cancel()
		err := srv.Run(ctx)
		assert.Nil(t, err)
		t.Logf("srv run on: %s", srv.addr)

		// wait ctx done and srv quit
		<-srv.Quit()
		assert.GreaterOrEqual(t, time.Since(st), runTime)
	})

	t.Run("test close by user", func(t *testing.T) {
		st := time.Now()

		// set run time
		runTime := time.Second * 3

		// run srv
		srv := Server{quit: make(chan struct{}, 1)}
		ctx, cancel := context.WithTimeout(context.Background(), runTime)
		defer cancel()
		err := srv.Run(ctx)
		assert.Nil(t, err)
		t.Logf("srv run on: %s", srv.addr)

		// wait for second and close by user
		time.Sleep(time.Second)
		srv.Close()

		// wait for srv quit
		<-srv.Quit()
		assert.Less(t, time.Since(st), runTime)
	})

	t.Run("test run error bind already used addr", func(t *testing.T) {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Error(err.Error())
			return
		}
		defer listener.Close()
		addr := listener.Addr().String()

		go func() {
			blockSrv := http.Server{Addr: addr}
			t.Logf("block srv already bind: %s", addr)
			go blockSrv.ListenAndServe()
			time.Sleep(time.Second * 2)
			_ = blockSrv.Shutdown(context.TODO())
		}()

		time.Sleep(time.Second)
		srv := Server{addr: addr, quit: make(chan struct{}, 1)}
		err = srv.Run(context.TODO())
		var targetErr *net.OpError
		assert.ErrorAs(t, err, &targetErr)
	})
}
