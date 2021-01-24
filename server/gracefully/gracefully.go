package gracefully

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(listenAndServe func() error, teardown func(context.Context) error) error {
	term := make(chan os.Signal)
	fail := make(chan error)
	go func() {
		signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
		<-term
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		fail <- teardown(ctx)
	}()
	if err := listenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return <-fail
}
