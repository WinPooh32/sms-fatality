package sysenv

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func CatchSignals(ctx context.Context, cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGSTOP,
	)

	select {
	case <-ctx.Done():
	case s := <-ch:
		log.Println("caught os signal:", s)
		log.Println("cancel tasks")

		cancel()

		signal.Stop(ch)
		close(ch)
	}
}
