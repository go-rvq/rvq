#!/bin/sh
# Builds the IDE into dist/ (gadide.js, gadide.css; embedded by ../gadide.go). GAD_DIR: a gad
# checkout with its submodules (web/ide-vuetify, web/plugins/js/*); by default
# the one beside this workspace.
set -e
cd "$(dirname "$0")"
bun install
bun run build
