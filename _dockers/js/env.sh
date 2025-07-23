#!/bin/bash
sourced=0

if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
  DIR=$(dirname "$0")
else
  DIR=$(dirname "${BASH_SOURCE[${#BASH_SOURCE[@]} - 1]}")
  sourced=1
fi

DIR=$(realpath "$DIR")

set -e

node_version=$(grep FROM "$DIR/Dockerfile" | perl -pe 's/^FROM node:(\d+)-.+$/\1/')
export NODE_CONTAINER_NAME="rvq_js:${node_version}-latest"

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

function run_container() {
  set -ex
  local opt="-it"

  case $1 in
  "-")
    opt="-i"
    ;;
  esac

  docker run "$opt" --rm \
        --user root \
        -e AS_USER=$(id -u):$(id -g) \
        -e HOME=/home/node \
        -v .:/home/node/src \
        -v "$DIR/volumes/pnpm-cache:/home/node/.cache/pnpm" \
        -v "$DIR/volumes/pnpm-store:/home/node/src/.pnpm-store" \
        -w /home/node/src \
        "$NODE_CONTAINER_NAME" "$@"

  local status=$?
  [ -e .pnpm-store ] && rmdir .pnpm-store
  return $status
}

function docker_pj_exec() {
  if [ $# -lt 2 ]; then
    echo "usage: pj_exec PJ_DIR CMD ..." >&2
    return 1
  fi

  cd "$1" || return 1

  shift

  run_container "$@" || return $?

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

  echo 'pnpm install && pnpm format && pnpm run build' | run_container -
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

eval "$@"