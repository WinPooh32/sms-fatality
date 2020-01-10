package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"common/broker"
	"common/fatality"
	"common/formats"
	"common/options"
	"common/sysenv"
)

const httpMaxHeaderBytes = 1 << 20

func serveInBackground(ctx context.Context, cancel context.CancelFunc, opts options.OptionsAcceptor) (*http.Server, error) {
	listener, err := net.ListenTCP (
		"tcp",
		&net.TCPAddr {
			IP:   net.ParseIP(opts.HttpIP),
			Port: int(opts.HttpPort),
			Zone: "",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("net listener: %w", err)
	}

	srv := &http.Server{
		Handler:        registerHandlers(http.NewServeMux()),
		ReadTimeout:    time.Duration(opts.HttpReadTimeout)  * time.Second,
		WriteTimeout:   time.Duration(opts.HttpWriteTimeout) * time.Second,
		MaxHeaderBytes: httpMaxHeaderBytes,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err = srv.Serve(listener)
		fatality.LogMsg("http serve: ", err)
	}()

	return srv, nil
}

// httpShutdownOnDone shutdowns http server on context done or server error
func httpShutdownOnDone(ctx, ctxServer context.Context, timeout uint32, srv *http.Server) {
	select {
	case <-ctx.Done():
		// gracefully shutdown a http server
		ctxShutdown, _ := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
		srv.Shutdown(ctxShutdown)

	case <-ctxServer.Done():
	}
}

func try( ctx context.Context, opts options.OptionsAcceptor, retry int) {
	const maxRecoveryRetries = 16

	if retry > maxRecoveryRetries {
		log.Println("recovery retries limit is reached")
		return
	}

	defer func() {
		if r := recover(); r!= nil {
			log.Println("acceptor recovered from ", r)
			log.Println("try again")

			try(ctx, opts, retry + 1)
		}
	}()

serviceLoop:
	for {
		ctxServer, cancelServer := context.WithCancel(ctx)

		publ, err := broker.ConnectAsPublisher(
			formats.AMQP(
				opts.MqHost,
				opts.MqPort,
				opts.MqUser,
				opts.MqPassword,
				opts.MqTimeout,
			),
		)
		if err != nil {
			log.Println("service failed to connect: ", err)
		} else {
			publ.HoldAlive(ctxServer)

			// put publisher connection object to context
			ctxWithData := context.WithValue(ctxServer, contextKeyPublisher, publ)
			//put server cancelation callback
			ctxWithData = context.WithValue(ctxWithData, contextKeyCancel, cancelServer)

			log.Printf("http server is listening on: %s:%d", opts.HttpIP, opts.HttpPort)
			srv, err := serveInBackground(ctxWithData, cancelServer, opts)
			if err != nil {
				log.Println("make new http server: ", err)
			} else {
				// wait for server grace shutdown complete
				httpShutdownOnDone(ctx, ctxServer, opts.HttpGracefulTimeout, srv)
			}

			publ.Close()
		}

		cancelServer()

		select {
		case <-ctx.Done():
			break serviceLoop
		default:
			log.Println("will restart in", opts.RestartDelay, "seconds")
			timeout := time.NewTimer(time.Duration(opts.RestartDelay) * time.Second)
			<-timeout.C
		}
	}
}

func main() {
	log.Println("starting acceptor")

	// define context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// stop service when os signal will be caught or context will be done
	go sysenv.CatchSignals(ctx, cancel)

	// read options from config file or app's arguments
	opts, err := options.ReadAcceptor()
	fatality.LogMsg("options.ReadWorker()", err)

	// try to run service
	try(ctx, opts, 0)

	log.Println("acceptor is shutting down")
}