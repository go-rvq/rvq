#!/bin/bash
sourced=0

if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
  DIR=$(dirname "$0")
else
  DIR=$(dirname "${BASH_SOURCE[${#BASH_SOURCE[@]} - 1]}")
  sourced=1
fi

[ "$NVM_DIR" = '' ] && {
  if [ -e ~/.nvm ]; then
    export NVM_DIR=~/.nvm
  else
    echo "NVM_DIR env is empty. See https://github.com/nvm-sh/nvm for install nvm instructions"

    [ $sourced -eq 1 ] && return  1 || exit 1
  fi
}

set -e

. "$NVM_DIR/nvm.sh"

node_version=$(grep FROM "$DIR/Dockerfile" | perl -pe 's/^FROM node:(\d+)-.+$/\1/')

nvm ls | grep $node_version >/dev/null || nvm install $node_version

nvm use $node_version >&2
export NODE_CONTAINER_NAME="node:${node_version}-latest-pnpm"

export PS1="(node `node --version`) \[\e]0;\u@\h: \w\a\]${debian_chroot:+($debian_chroot)}\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ "

function docker_image_build() {
  ( cd "$DIR" && \
    docker build --tag $NODE_CONTAINER_NAME .
  )
  return $?
}

function docker_image_rm() {
  docker image rm $NODE_CONTAINER_NAME
  return $?
}

function docker_image_check() {
  echo INFO: checking build image >&2
  docker inspect --type image $NODE_CONTAINER_NAME 2>/dev/null >&2 || \
    docker_image_build >&2
  return $?
}

docker_image_check

function build_pj() {
  if [ $# -eq 0 ]; then
    echo "usage: build_pj PJ_DIR" >&2
    return 1
  fi

  cd "$1" || return 1
  pnpm install && pnpm format && pnpm run build
  return $?
}

function docker_pj_exec() {
  if [ $# -lt 2 ]; then
    echo "usage: pj_exec PJ_DIR CMD ..." >&2
    return 1
  fi

  cd "$1" || return 1

  shift

  docker run -it --rm \
      --user $(id -u):$(id -g) \
      -e HOME=/home/node \
      -v .:/home/node/src \
      -v ~/.cache/pnpm:/home/node/.cache/pnpm \
      -v $(dirname $(pnpm store path)):/home/node/src/.pnpm-store \
      -w /home/node/src \
      "$NODE_CONTAINER_NAME" "$@"

  "$@"
  return $?
}

function docker_build_pj() {
  if [ $# -eq 0 ]; then
    echo "usage: docker_build_pj PJ_DIR" >&2
    return 1
  fi

  cd "$1" || return 1
  [ -e ./dist ] && rm -vrf ./dist

  docker run -it --rm \
    --user $(id -u):$(id -g) \
    -e HOME=/home/node \
    -v .:/home/node/src \
    -v ~/.cache/pnpm:/home/node/.cache/pnpm \
    -v $(dirname $(pnpm store path)):/home/node/src/.pnpm-store \
    -w /home/node/src \
    pnpm:latest bash -c 'pnpm install && pnpm format && pnpm run build' || return $?

  [ -e .pnpm-store ] && rmdir .pnpm-store
  return $?
}

function pj_exec() {
  if [ $# -lt 2 ]; then
    echo "usage: pj_exec PJ_DIR CMD ..." >&2
    return 1
  fi

  cd "$1" || return 1

  shift

  "$@"
  return $?
}

function pj_pnpm() {
  if [ $# -lt 2 ]; then
    echo "usage: pj_pnpm PJ_DIR PNPM_ARGS..." >&2
    return 1
  fi

  local dir="$1"

  shift

  pj_exec "$dir" pnpm install || return $?
  pj_exec "$dir" pnpm "$@"
  return $?
}

if [ $sourced -eq 0 ]; then
  if [ $# -eq 0 ]; then
    exec bash --norc
  else
    "$@"
  fi
fi