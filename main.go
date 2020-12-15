package main

import (
	"net/http"
	"os"

	"github.com/findy-network/findy-agent-vault/resolver"
	"github.com/findy-network/findy-agent-vault/utils"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-vault/server"
)

const defaultPort = "8085"

var gqlResolver *resolver.Resolver

func main() {
	utils.SetLogDefaults()
	gqlResolver = resolver.InitResolver(false)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	srv := server.Server(gqlResolver)
	http.Handle("/query", srv)

	glog.Fatal(http.ListenAndServe(":"+port, nil))
}
