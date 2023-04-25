package command

import (
	"challenge/api"
	"challenge/api/actor"
	"challenge/config"
	"challenge/model"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
	"net/http"
)

var serveHandler cli.ActionFunc = func(c *cli.Context) error {
	rtr := mux.NewRouter()
	rtr.Use(jsonableMiddleware)

	models := c.Value(model.ApplicationModelsContext).(model.ApplicationModels)

	http.Handle("/", rtr)

	rtr.HandleFunc("/healthcheck", api.Healthcheck).Methods(http.MethodGet)

	rtr.HandleFunc("/actor", actor.AddActor(&models)).Methods(http.MethodPost)
	rtr.HandleFunc("/actor", actor.UpdateActor(&models)).Methods(http.MethodPut)
	rtr.HandleFunc("/actor", actor.GetActor(&models)).Methods(http.MethodGet)
	rtr.HandleFunc("/actors", actor.GetActors(&models)).Methods(http.MethodGet)

	log.Infof("Listening on %s", config.Cfg().Listen)

	server := &http.Server{
		Addr:    config.Cfg().Listen,
		Handler: nil,
		BaseContext: func(l net.Listener) context.Context {
			return c.Context
		},
	}

	go func() {
		<-c.Done()
		_ = server.Shutdown(c.Context)
	}()

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func jsonableMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

var ServeCommand = &cli.Command{
	Name:        "serve",
	Description: "Start a HTTP(s) server",
	Usage:       "Start a HTTP(s) server",
	Action:      serveHandler,
}

var ServerCommands = &cli.Command{
	Name:        "server",
	Description: "HTTP(s) server commands",
	Usage:       "HTTP(s) server commands",
	Subcommands: []*cli.Command{
		ServeCommand,
	},
}
