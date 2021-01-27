BRANCH:=dev

clean:
	docker-compose down
	rm -rf .docker

clone:
	mkdir -p .docker
	-git clone git@github.com:findy-network/findy-agent-cli.git .docker/findy-agent-cli
	cd .docker/findy-agent-cli && git checkout $(BRANCH)
	-git clone git@github.com:findy-network/findy-grpc.git .docker/findy-grpc
	cd .docker/findy-grpc && git checkout $(BRANCH)

clean_cli:
	docker rmi findy-agent-cli

build_cli:
	cd .docker/findy-agent-cli && make image

up: clone
	docker-compose up

prepare-test:
	-docker-compose down
	docker-compose up --build -d
	sleep 35
	docker build -t vault-demo-test --build-arg HTTPS_PREFIX=$(HTTPS_PREFIX) ./test

run-test:
	docker run --name vault-demo-e2e \
		--rm -it --network=demo_default \
		-v $(PWD)/.docker/findy-grpc/cert:/go/src/github.com/findy-network/findy-grpc/cert \
		-v $(PWD)/test/test.sh:/test.sh \
		vault-demo-test

test: prepare-test run-test

tree:
	docker run --rm -it findy-agent-cli tree