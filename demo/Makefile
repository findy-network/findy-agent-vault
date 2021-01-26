BRANCH:=dev

clone:
	mkdir -p .docker
	-git clone git@github.com:findy-network/findy-agent-cli.git .docker/findy-agent-cli
	cd .docker/findy-agent-cli && git checkout $(BRANCH)
	-git clone git@github.com:findy-network/findy-grpc.git .docker/findy-grpc
	cd .docker/findy-grpc && git checkout $(BRANCH)

build_cli:
	cd .docker/findy-agent-cli && make image

up: clone
	docker-compose up

run-test:
	docker-compose up -d
	sleep 20
	docker build -t vault-demo-test --build-arg HTTPS_PREFIX=$(HTTPS_PREFIX) ./test
	docker run --name vault-demo-e2e \
		--rm -it --network=demo_default \
		-v $(PWD)/.docker/findy-grpc/cert:/go/src/github.com/findy-network/findy-grpc/cert \
		-v $(PWD)/test/test.sh:/test.sh \
		vault-demo-test

tree:
	docker run --rm -it findy-agent-cli tree