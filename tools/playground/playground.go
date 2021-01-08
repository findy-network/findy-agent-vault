package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/lainio/err2"

	"github.com/rs/cors"

	"github.com/golang/glog"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/findy-network/findy-agent-vault/resolver"
	"github.com/findy-network/findy-agent-vault/server"
	"github.com/findy-network/findy-agent-vault/utils"
)

var srv *server.VaultServer

func TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			glog.Error("ERROR generating token:", err.Error())
		})

		token, err := srv.CreateToken(uuid.New().String())
		err2.Check(err)

		w.Header().Add("Content-Type", "text/plain")
		_, err = w.Write([]byte(token))
		err2.Check(err)
	}
}

func main() {
	utils.SetLogDefaults()
	config := utils.LoadConfig()
	srv = server.NewServer(resolver.InitResolver(true, true, true), config.JWTKey)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv.Handle())
	http.Handle("/token", cors.AllowAll().Handler(TokenHandler()))

	glog.Infof("connect to http://localhost:%d/ for GraphQL playground", config.ServerPort)
	log.Fatal(http.ListenAndServe(config.Address, nil))
}
