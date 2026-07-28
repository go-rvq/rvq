function exampleRestart() {
  echo "=================>"
  killall rvqexample
  source example/dev_env
#  export DEV_PRESETS=1
  go build -o /tmp/rvqexample example/main.go && /tmp/rvqexample
}

export -f exampleRestart

find . -name "*.go" | entr -r bash -c "exampleRestart"
