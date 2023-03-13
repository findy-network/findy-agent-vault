package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	defer gqlResolver.Close()

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
	startServer(config.Address)
}

func startServer(address string) {
	const serverTimeout = 5 * time.Second
	ourServer := &http.Server{
		Addr:              address,
		ReadHeaderTimeout: serverTimeout,
	}

	go func() {
		if err := ourServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			glog.Errorf("HTTP server error: %v", err)
		} else {
			glog.Infoln("Stopped serving new connections.")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), serverTimeout)
	defer shutdownRelease()

	if err := ourServer.Shutdown(shutdownCtx); err != nil {
		glog.Errorf("HTTP shutdown error: %v", err)
	} else {
		glog.Infoln("Graceful shutdown complete.")
	}
}
