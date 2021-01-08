package main

import (
	"net/http"

	"github.com/findy-network/findy-agent-vault/resolver"
	"github.com/findy-network/findy-agent-vault/utils"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-vault/server"
)

func main() {
	utils.SetLogDefaults()
	config := utils.LoadConfig()

	gqlResolver := resolver.InitResolver(false, false, false)

	srv := server.NewServer(gqlResolver, config.JWTKey)
	http.Handle("/query", srv.Handle())

	glog.Fatal(http.ListenAndServe(config.Address, nil))
}
