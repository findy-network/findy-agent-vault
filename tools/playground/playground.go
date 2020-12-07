package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lainio/err2"

	"github.com/rs/cors"

	"github.com/golang/glog"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/findy-network/findy-agent-vault/server"
	"github.com/findy-network/findy-agent-vault/tools/resolver"
	"github.com/findy-network/findy-agent-vault/utils"
)

const defaultPort = "8085"

var gqlResolver *resolver.Resolver

func TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			fmt.Println("ERROR generation token:", err.Error())
		})
		user, err := gqlResolver.Query().User(context.TODO())
		err2.Check(err)

		token, err := server.CreateToken(user.ID)
		err2.Check(err)

		w.Header().Add("Content-Type", "text/plain")
		_, err = w.Write([]byte(token))
		err2.Check(err)
	}
}

func main() {
	utils.SetLogDefaults()
	gqlResolver = resolver.InitResolver(false)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	srv := server.Server(gqlResolver)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.Handle("/token", cors.AllowAll().Handler(TokenHandler()))

	glog.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
