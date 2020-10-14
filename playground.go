package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/findy-network/findy-agent-api/tools/data"

	"github.com/golang/glog"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/findy-network/findy-agent-api/server"
	"github.com/findy-network/findy-agent-api/tools/resolver"
)

const defaultPort = "8085"

func initLogging() {
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")
	flag.Set("v", "3")
	flag.Parse()
}

func TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := server.CreateToken(data.State.User.ID)
		if err == nil {
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte(token))
		} else {
			panic(err)
		}
	}
}

func main() {
	initLogging()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	/*// TEST subscription
	ticker := time.NewTicker(time.Second * 30)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				resolver.AddEvent()
			}
		}
	}()
	// TEST SUBSCRIPTION end*/

	srv := server.Server(&resolver.Resolver{})
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.Handle("/token", cors.AllowAll().Handler(TokenHandler()))

	glog.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
