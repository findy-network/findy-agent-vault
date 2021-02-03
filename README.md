# findy-agent-vault

Data storage service for findy-agency clients. Service provides GraphQL interface for interaction.

## Running with mock data

Service providing mock data can be launched with following steps:

1. [Install go](https://golang.org/dl/)
2. Run
    
    ```
    make run
    ```

This will launch the service in port 8085.
* Access graphiQL playground with browser: http://localhost:8085
* Configure URL `http://localhost:8085/query` to your own GQL-client.

See [sample client implementation](https://github.com/findy-network/findy-wallet-pwa).

### Authentication

Running the service in playground mode provides an endpoint for mock authentication token generation.
Visit http://localhost:8085/token to generate the token.
API requests should contain this token in header field for the authentication to succeed:
```
{"Authorization": "Bearer <TOKEN>"}
```

## Running with postgres and findy-agent

1. Start postgres and findy-agent in their own docker containers:
    ```bash
    make dev_build
    ```

1. Onboard your agent to agency:

    ```bash
    go run tools/onboard/main.go
    ```
    Copy JWT token from the produced output.

1. Declare following environment variables:

    ```bash
    export FAV_SERVER_PORT=8085
    export FAV_USE_PLAYGROUND=true
    export FAV_AGENCY_PORT=50052
    export FAV_DB_PASSWORD="my-secret-password"
    export FAV_AGENCY_CERT_PATH=".github/workflows/cert"
    ```

1. Run vault:

    ```bash
    go run main.go
    ```
1. Open http://localhost:8085 and set token as [instructed](#authentication). Execute graphQL queries with the playground e.g.

    ```
    { user { id name } }
    ```

## Unit testing

Unit tests assume postgres is running on port 5433.

Launch default postgres container by declaring password for postgres user:
```bash
export FAV_DB_PASSWORD="mysecretpassword"
```

and running

```
make init-test
```

You can run all unit tests with command

```bash
go test ./...
````

For linting, you need to install [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

`make check` builds, tests and lints the code.

