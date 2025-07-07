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
		pnpm:latest bash -c 'pnpm install && pnpm run build'

	cp vue/* dist

build-format: docker_build
	$(DOCKER_CMD) run -it --rm \
		--user $(shell id -u):$(shell id -g) \
		-e HOME=/home/node \
		-v $(HOME)/.cache/pnpm:/home/node/.cache/pnpm \
		-v .:/home/node/src \
		-w /home/node/src \
		pnpm:latest bash -c 'pnpm install && pnpm format && pnpm run build'

	cp vue/* dist

get-vue:
	wget https://unpkg.com/vue@3/dist/vue.global.js -O vue/vue.global.dev.js
	wget https://unpkg.com/vue@3/dist/vue.global.prod.js -O vue/vue.global.prod.js