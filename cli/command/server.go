package command

import (
	"challenge/api"
	"challenge/config"
	"challenge/model"
	cjwt "challenge/service/jwt"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net"
	"net/http"
	"strings"
)

var serveHandler cli.ActionFunc = func(c *cli.Context) error {
	rtr := mux.NewRouter()
	rtr.Use(jsonableMiddleware)
	rtr.Use(authMiddleware)

	models := c.Context.Value(model.ApplicationModelsContext).(*model.ApplicationModels)

	http.Handle("/", rtr)

	rtr.HandleFunc("/healthcheck", api.Healthcheck).Methods(http.MethodGet)
	rtr.HandleFunc("/signin", api.SignIn).Methods(http.MethodPost)

	rtr.HandleFunc("/actor", api.AddActor(models)).Methods(http.MethodPost)
	rtr.HandleFunc("/actor/{uuid}", api.UpdateActor(models)).Methods(http.MethodPut)
	rtr.HandleFunc("/actor/{uuid}", api.GetActor(models)).Methods(http.MethodGet)
	rtr.HandleFunc("/actors", api.GetActors(models)).Methods(http.MethodGet)

	rtr.HandleFunc("/protected/hi", api.ProtectedResource).Methods(http.MethodGet)

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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/protected") { // not the best solution here. It was a last minute idea))
			next.ServeHTTP(w, r)

			return
		}

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			_, _ = api.CreateResponseUnauthorized(w, "Missing Authorization Header")

			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := cjwt.VerifyToken(tokenString)
		if err != nil {
			_, _ = api.CreateResponseUnauthorized(w, "Error verifying JWT token: "+err.Error())

			return
		}
		name := claims.(jwt.MapClaims)["name"].(string)
		role := claims.(jwt.MapClaims)["role"].(string)

		r.Header.Set("X-Name", name)
		r.Header.Set("X-Role", role)

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
