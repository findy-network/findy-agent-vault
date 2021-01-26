BRANCH:=dev

clone:
	mkdir -p .docker
	-git clone git@github.com:findy-network/findy-agent-cli.git .docker/findy-agent-cli
	cd .docker/findy-agent-cli && git checkout $(BRANCH)
	-git clone git@github.com:findy-network/findy-grpc.git .docker/findy-grpc
#	cd .docker/findy-grpc && git checkout $(BRANCH)
	cd .docker/findy-grpc && git checkout add-dockerfile

build_cli:
	cd .docker/findy-agent-cli && make image

up: clone build_cli
	docker-compose up