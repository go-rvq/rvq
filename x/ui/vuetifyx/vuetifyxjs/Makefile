third_party_dir := $(shell realpath "../_third_party/js")
prismjs_dir=$(shell realpath $(third_party_dir)"/prismjs")

docker_args="docker_args.start docker_args.end"

# ../dev is a local repository of https://github.com/go-rvq/dev
docker_image_env = ../../../../../dev/dockers/js/env.sh

.PHONY: dev
dev:
	rvq_js_container_name=vuetifyxjs_dev $(docker_image_env) @dev "$(docker_args)"

.PHONY: install
install:
	echo 'pnpm install' | $(docker_image_env) @exec - "$(docker_args)"

.PHONY: bash
bash:
	$(docker_image_env) @exec "$(docker_args)" bash

.PHONY: build
build:
	$(docker_image_env) @build "$(docker_args)"

build_image:
	$(docker_image_env) @image_build


.PHONY: dist-clean
dist-clean:
	rm -rf dist

.PHONY: clean
clean:
	rm -rf node_modules .pnpm-store

teste:
	$(docker_image_env) @exec "$(docker_args)" ls -la .docker/mount-points/codemirror-gad
