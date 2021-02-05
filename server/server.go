package server

import (
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/findy-network/findy-agent-vault/graph/generated"
	"github.com/findy-network/findy-agent-vault/utils"
	jwtMW "github.com/findy-network/findy-grpc/jwt/mw"
	"github.com/gorilla/websocket"
)

const (
	queryCacheSize          = 1000
	persistedQueryCacheSize = 100
)

type VaultServer struct {
	server      *handler.Server
	authChecker *jwtMW.Middleware
}

func schema(resolver generated.ResolverRoot) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.LogTrace().Infof("received request: %s %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func NewServer(resolver generated.ResolverRoot, jwtSecret string) *VaultServer {
	srv := handler.New(schema(resolver))

	// TODO: figure out CORS policy for our WS use case
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: true,
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(queryCacheSize))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(persistedQueryCacheSize),
	})

	authChecker := jwtMW.New(&jwtMW.Options{
		ValidationKey:  jwtSecret,
		QueryParameter: "access_token",
		ErrorHandler:   onAuthError,
	})

	return &VaultServer{
		server:      srv,
		authChecker: authChecker,
	}
}

func (v *VaultServer) Handle() http.Handler {
	// TODO: figure out CORS policy for our HTTP use case
	return cors.AllowAll().Handler(logRequest(v.authChecker.Handler(v.server)))
}
