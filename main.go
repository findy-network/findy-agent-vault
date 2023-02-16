package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/findy-network/findy-agent-vault/agency/findy"
	"github.com/findy-network/findy-agent-vault/resolver"
	"github.com/findy-network/findy-agent-vault/server"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang/glog"
)

func main() {
	utils.SetLogDefaults()
	config := utils.LoadConfig()

	if len(os.Args) > 1 && os.Args[1] == "version" {
		log.Printf("Vault version %s\n", config.Version)
		return
	}

	gqlResolver := resolver.InitResolver(config, &findy.Agency{})

	srv := server.NewServer(gqlResolver, config.JWTKey)
	http.Handle("/query", srv.Handle())
	if config.UsePlayground {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if utils.LogTrace() {
			glog.Infof("health check %s %s", r.URL.Path, config.Version)
		}
		_, _ = w.Write([]byte(config.Version))
	})

	const serverTimeout = 5 * time.Second
	ourServer := &http.Server{
		Addr:              config.Address,
		ReadHeaderTimeout: serverTimeout,
	}

	glog.Fatal(ourServer.ListenAndServe())
}
