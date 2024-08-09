DOCKER_CMD ?= docker

docker_build:
	$(DOCKER_CMD) build --tag pnpm:latest .

build: docker_build
	$(DOCKER_CMD) run -it --rm \
		--user $(shell id -u):$(shell id -g) \
		-e HOME=/home/node \
		-v $(HOME)/.cache/pnpm:/home/node/.cache/pnpm \
		-v .:/home/node/src \
		-w /home/node/src \
		pnpm:latest bash -c 'pnpm install && pnpm format && pnpm run build'