package main

import (
	"challenge/cli"
	"challenge/config"
	"challenge/model"
	"context"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const TimezoneUTC = "UTC"

func main() {
	var err error
	var cfg = config.Cfg()
	loc, _ := time.LoadLocation(TimezoneUTC)
	time.Local = loc

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"mysql://%s:%s@%s:%s/%s",
			cfg.DBConfig.User,
			cfg.DBConfig.Password,
			cfg.DBConfig.Host,
			cfg.DBConfig.Port,
			cfg.DBConfig.Name,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	models := &model.ApplicationModels{
		Actors: model.ActorModel{DB: db},
	}
	ctx, cancel := context.WithCancel(
		context.WithValue(
			context.Background(),
			model.ApplicationModelsContext,
			models,
		),
	)

	done := make(chan bool)
	go func() {
		defer close(done)
		err = cli.NewApplication().RunContext(ctx, os.Args)
		if err != nil {
			log.Fatalf("error running application: %v", err)
		} else {
			log.Info("Graceful shutdown")
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-wait: // wait for SIGINT/SIGTERM
		signal.Reset(syscall.SIGINT, syscall.SIGTERM) // resetting signal listener, so that repeated Ctrl+C will exit immediately
		cancel()                                      // graceful stop
		<-done

	case <-done:
		cancel() // graceful stop
	}
}
