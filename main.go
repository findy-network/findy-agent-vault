package main

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/findy-network/findy-agent-vault/resolver"
	"github.com/findy-network/findy-agent-vault/utils"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-vault/server"
)

func main() {
	utils.SetLogDefaults()
	config := utils.LoadConfig()

	gqlResolver := resolver.InitResolver(config)

	srv := server.NewServer(gqlResolver, config.JWTKey)
	http.Handle("/query", srv.Handle())
	if config.UsePlayground {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	glog.Fatal(http.ListenAndServe(config.Address, nil))
}
