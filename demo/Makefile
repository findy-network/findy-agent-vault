BRANCH:=dev
TIMEOUT?=35

clean:
	docker-compose -f docker-compose-test.yml down 
	docker-compose -f docker-compose-dev.yml down 
	rm -rf .docker

clone:
	mkdir -p .docker
	-git clone git@github.com:findy-network/findy-agent.git .docker/findy-agent
	cd .docker/findy-agent && git checkout $(BRANCH)
	-git clone git@github.com:findy-network/findy-common-go.git .docker/findy-common-go
	cd .docker/findy-common-go && git checkout $(BRANCH)

clean_cli:
	docker rmi findy-agent

build_cli:
	cd .docker/findy-agent && make image

prepare-test:
	-docker-compose -f docker-compose-test.yml down
	docker-compose -f docker-compose-test.yml up --build -d
	mkdir -p ./test/.docker
	cp -R .docker/findy-common-go/cert ./test/.docker/cert
	docker build -t vault-demo-test --build-arg HTTPS_PREFIX=$(HTTPS_PREFIX) ./test
	sleep $(TIMEOUT)

run-test:
	docker run --name vault-demo-e2e \
		--rm -i --network=demo_default \
		vault-demo-test

test: prepare-test run-test

tree:
	docker run --rm -it findy-agent tree

test-test:
	docker run --name vault-demo-e2e \
		--rm -it --network=demo_default \
		-v $(PWD)/.docker/findy-common-go/cert:/go/src/github.com/findy-network/findy-common-go/cert \
		-v $(PWD)/test/test.sh:/test.sh \
		vault-demo-test

dev:
	docker-compose -f docker-compose-dev.yml down
	rm -rf .docker/db
	docker-compose -f docker-compose-dev.yml up -d

dev_build: clean clean_cli clone build_cli
	docker-compose -f docker-compose-dev.yml up --build -d

dev_logs:
	docker-compose -f docker-compose-dev.yml logs --tail="all"
