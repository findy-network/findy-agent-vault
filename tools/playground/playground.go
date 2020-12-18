package main

import (
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/lainio/err2"

	"github.com/rs/cors"

	"github.com/golang/glog"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/findy-network/findy-agent-vault/resolver"
	"github.com/findy-network/findy-agent-vault/server"
	"github.com/findy-network/findy-agent-vault/utils"
)

const defaultPort = "8085"

func TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			glog.Error("ERROR generation token:", err.Error())
		})

		// TODO
		token, err := server.CreateToken(uuid.New().String())
		err2.Check(err)

		w.Header().Add("Content-Type", "text/plain")
		_, err = w.Write([]byte(token))
		err2.Check(err)
	}
}

func main() {
	utils.SetLogDefaults()
	srv := server.Server(resolver.InitResolver(true, true))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.Handle("/token", cors.AllowAll().Handler(TokenHandler()))

	glog.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
