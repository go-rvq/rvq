DOCKER_CMD ?= docker

docker_args="docker_args.start docker_args.end"

# ../dev is a local repository of https://github.com/go-rvq/dev
docker_image_env = ../../../dev/dockers/js/env.sh

.PHONY: build
build:
	$(docker_image_env) @build "$(docker_args)"
	cp assets/* dist

get:
	wget https://unpkg.com/vue@3/dist/vue.global.js -O assets/vue.global.dev.js
	wget https://unpkg.com/vue@3/dist/vue.global.prod.js -O assets/vue.global.prod.js
	wget https://unpkg.com/vue-i18n@11 -O assets/vue-i18n.js