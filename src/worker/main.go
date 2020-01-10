package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"time"

	"common/broker"
	"common/db"
	"common/fatality"
	"common/formats"
	"common/options"
	"common/sysenv"

	_ "github.com/lib/pq" // postgresql driver
)

const (
	// will be compared with database version on start
	dbVersion = 2
)

func validateVersion(ctx context.Context, dbConn *sql.DB) error {
	dbQueries := db.New(dbConn)

	versionRow, err := dbQueries.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("validate version: %w", err)
	}

	if versionRow.VersionID != dbVersion || !versionRow.IsApplied {
		return fmt.Errorf("worker and database versions are not equal: %d", versionRow.VersionID)
	}

	return nil
}

func connect(opts options.OptionsWorker) (*sql.DB, broker.Consumer, error){
	log.Printf("connecting to database at %s:%d", opts.DbHost, opts.DbPort)
	pg, err := sql.Open("postgres",
		formats.Postgres(
			opts.DbHost,
			opts.DbPort,
			opts.DbName,
			opts.DbUser,
			opts.DbPassword,
			opts.DbSsl,
			opts.DbTimeout,
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("database: %w", err)
	}

	log.Printf("connecting to rabbitMQ at %s:%d", opts.MqHost, opts.MqPort)
	consumer, err := broker.ConnectAsConsumer(
		formats.AMQP(
			opts.MqHost,
			opts.MqPort,
			opts.MqUser,
			opts.MqPassword,
			opts.MqTimeout,
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("rabbitMQ broker: %w", err)
	}

	return pg, consumer, nil
}

func init(){
	// service should be binded to one cpu core for better performance
	// 1 service per core
	runtime.GOMAXPROCS(1)
}

func main(){
	ctx, cancel := context.WithCancel(context.Background())

	go sysenv.CatchSignals(ctx, cancel)

	// read options from config file or app's arguments
	opts, err := options.ReadWorker()
	if err != nil{
		log.Fatalln("options.ReadWorker(): ", err)
	}

serviceLoop:
	for {
		connPg, consumer, err := connect(opts)
		if err != nil {
			log.Println("service failed to connect: ", err)
		} else {
			err = validateVersion(ctx, connPg)
			fatality.Panic(err);

			log.Println("started to process messages")
			work(ctx, connPg, consumer)

			err := consumer.Close()
			fatality.Log(err)

			err = connPg.Close()
			fatality.Log(err)
		}

		select {
		case <-ctx.Done():
			break serviceLoop
		default:
			log.Println("will restart in", opts.RestartDelay, "seconds")
			timeout := time.NewTimer(time.Duration(opts.RestartDelay) * time.Second)
			<-timeout.C
		}
	}

	log.Println("worker is shutting down")
}