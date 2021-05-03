S3_TOOL_ASSET_PATH := https://$(HTTPS_PREFIX)api.github.com/repos/findy-network/findy-common-go/releases/assets/34026533

.PHONY: db

scan:
	@./scan.sh

run:
	go run tools/playground/playground.go

generate: 
	go generate ./...

deps:
	go get -t ./...

update-deps:
	go get -u ./...

build: generate
	go build -v ./...

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
		-e POSTGRES_PASSWORD=$(FAV_DB_PASSWORD) \
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
		-e POSTGRES_PASSWORD=$(FAV_DB_PASSWORD) \
		-e POSTGRES_DB=vault \
		-v $(PWD)/.db/data:/var/lib/postgresql/data \
		-p 5432:5432 \
		-d postgres:13.1-alpine

db_client:
	docker run -it --rm --network host jbergknoff/postgresql-client postgres://postgres:$(FAV_DB_PASSWORD)@localhost:5432/vault?sslmode=disable

db_client_test:
	docker run -it --rm --network host jbergknoff/postgresql-client postgres://postgres:$(FAV_DB_PASSWORD)@localhost:5433/vault?sslmode=disable

check:
	go build ./...
	go test ./...
	golangci-lint run

remod:
	-rm go*
	go mod init github.com/findy-network/findy-agent-vault
	go mod tidy

single_test:
	go test -run TestConnect ./...

dclean:
	-docker rmi findy-agent-vault

dbuild:
	@[ "${HTTPS_PREFIX}" ] || ( echo "ERROR: HTTPS_PREFIX <{githubUser}:{githubToken}@> is not set"; exit 1 )
	mkdir -p .docker
	curl -L -H "Accept:application/octet-stream" "$(S3_TOOL_ASSET_PATH)" > .docker/s3-copy
	docker build \
		--build-arg HTTPS_PREFIX=$(HTTPS_PREFIX) \
		-t findy-agent-vault \
		.

gen_mock:
	~/go/bin/mockgen -package listen -source ./db/store/db.go DB > ./resolver/listen/listener_mock_store_test.go
