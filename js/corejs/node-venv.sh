#!/bin/bash
sourced=0

if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
  DIR=$(dirname "$0")
else
  DIR=$(dirname "${BASH_SOURCE[${#BASH_SOURCE[@]} - 1]}")
  sourced=1
fi

[ "$NVM_DIR" = '' ] && {
  echo "NVM_DIR env is empty. See https://github.com/nvm-sh/nvm for install nvm instructions"
  return 1
}

set -e

. "$NVM_DIR/nvm.sh"

node_version=$(grep FROM "$DIR/Dockerfile" | perl -pe 's/^FROM node:(\d+)-.+$/\1/')

nvm ls | grep $node_version >/dev/null || nvm install $node_version

nvm use $node_version

(
  cd "$DIR" && \
  pnpm install
)

export PS1="(node `node --version`) \[\e]0;\u@\h: \w\a\]${debian_chroot:+($debian_chroot)}\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ "

if [ $sourced -eq 0 ]; then
  if [ $# -eq 0 ]; then
    bash --norc
  fi
fi