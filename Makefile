.PHONY: db

scan:
	@./scripts/scan.sh $(ARGS)

scan_and_report:
	@./scripts/scan.sh v > licenses.txt

generate: 
	go generate ./...

deps:
	go get -t ./...

update-deps:
	go get -u ./...

build:
	go build -v ./...

shadow:
	@echo Running govet
	go vet -vettool=$(GOPATH)/bin/shadow ./...
	@echo Govet success

check_fmt:
	$(eval GOFILES = $(shell find . -name '*.go'))
	@gofmt -l $(GOFILES)

lint:
	golangci-lint run

init-test:
	-docker stop findy-agent-vault-test-db
	-docker rm findy-agent-vault-test-db
	-rm -rf .db/test
	docker run --name findy-agent-vault-test-db \
		-e POSTGRES_PASSWORD=$(FAV_DB_PASSWORD) \
		-e POSTGRES_DB=vault \
		-v $(PWD)/.db/test:/var/lib/postgresql/data \
		-p 5433:5432 \
		-d postgres:13.8-alpine
	sleep 30


test:
	go test -v ./...

test_cov_out:
	go test \
		-coverpkg=github.com/findy-network/findy-agent-vault/... \
		-coverprofile=coverage.txt  \
		-covermode=atomic \
		./...

test_cov: test_cov_out
	go tool cover -html=coverage.txt

db:
	-docker stop findy-agent-vault-db
	-docker rm findy-agent-vault-db
	-rm -rf .db/data
	docker run --name findy-agent-vault-db \
		-e POSTGRES_PASSWORD=$(FAV_DB_PASSWORD) \
		-e POSTGRES_DB=vault \
		-v $(PWD)/.db/data:/var/lib/postgresql/data \
		-p 5432:5432 \
		-d postgres:13.8-alpine

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
	docker build -t findy-agent-vault .

gen_mock:
	go install github.com/golang/mock/mockgen@v1.6.0
	~/go/bin/mockgen -package listen -source ./db/store/db.go DB > ./resolver/listen/listener_mock_store_test.go
	~/go/bin/mockgen -package mock -source ./agency/model/model.go Agency > ./agency/mock/mock.go
