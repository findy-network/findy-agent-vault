# findy-agent-vault

Data storage service for findy-agency clients. Service provides GraphQL interface for interaction.

## Usage

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

## Authentication

Running the service in playground mode provides an endpoint for mock authentication token generation.
Visit http://localhost:8085/token to generate the token.
API requests should contain this token in header field for the authentication to succeed:
```
{"Authorization": "Bearer <TOKEN>"}
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

