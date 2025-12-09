DOCKER_CMD ?= docker

build:
	../../_dockers/js/env.sh @build
	cp assets/* dist

get:
	wget https://unpkg.com/vue@3/dist/vue.global.js -O assets/vue.global.dev.js
	wget https://unpkg.com/vue@3/dist/vue.global.prod.js -O assets/vue.global.prod.js
	wget https://unpkg.com/vue-i18n@11 -O assets/vue-i18n.js