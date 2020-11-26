
run:
	go run playground.go

generate: 
	go generate ./...

deps:
	go get -t ./...

update-deps:
	go get -u ./...

build: generate
	go build -v ./...

build_findy: generate
	go build -tags findy -v ./...

build_findy_grpc: generate
	go build -tags findy_grpc -v ./...

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

test:
	go test -v ./...

test_cov:
	go test -v -coverprofile=c.out ./... && go tool cover -html=c.out

check: check_fmt vet shadow

