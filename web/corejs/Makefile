DOCKER_CMD ?= docker

build:
	../../_dockers/js/env.sh docker_build_pj $(shell pwd)
	cp vue/* dist

get-vue:
	wget https://unpkg.com/vue@3/dist/vue.global.js -O vue/vue.global.dev.js
	wget https://unpkg.com/vue@3/dist/vue.global.prod.js -O vue/vue.global.prod.js