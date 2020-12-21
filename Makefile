.PHONY: db

run:
	go run tools/playground/playground.go

generate: 
	go generate ./...

deps:
	go get -t ./...

update-deps:
	go get -u ./...

build_findy: generate
	go build -tags findy -v ./...

build: generate
	go build -v ./...

test_findy:
	go test -tags findy -v ./...

vet:
	go vet ./...

shadow:
	@echo Running govet
	go vet -vettool=$(GOPATH)/bin/shadow ./...
	@echo Govet success

check_fmt:
	$(eval GOFILES = $(shell find . -name '*.go'))
	@gofmt -l $(GOFILES)

lint:
	$(GOPATH)/bin/golint ./... 

lint_e:
	@$(GOPATH)/bin/golint ./... | grep -v export | cat

init-test:
	-docker stop findy-agent-vault-test-db
	-docker rm findy-agent-vault-test-db
	-rm -rf .db/test
	docker run --name findy-agent-vault-test-db \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=vault \
		-v $(PWD)/.db/test:/var/lib/postgresql/data \
		-p 5433:5432 \
		-d postgres:13.1-alpine
	sleep 30


test:
	go test -v ./...

test_cov:
	go test -coverprofile=c.out ./... -coverpkg=./... && go tool cover -html=c.out

db:
	-docker stop findy-agent-vault-db
	-docker rm findy-agent-vault-db
	-rm -rf .db/data
	docker run --name findy-agent-vault-db \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=vault \
		-v $(PWD)/.db/data:/var/lib/postgresql/data \
		-p 5432:5432 \
		-d postgres:13.1-alpine

db_client:
	docker run -it --rm --network host jbergknoff/postgresql-client postgres://postgres:$(POSTGRES_PASSWORD)@localhost:5432/vault?sslmode=disable

db_client_test:
	docker run -it --rm --network host jbergknoff/postgresql-client postgres://postgres:$(POSTGRES_PASSWORD)@localhost:5433/vault?sslmode=disable

check:
	go build -tags findy_grpc ./...
	go test -tags findy_grpc ./...
	golangci-lint --build-tags findy_grpc run