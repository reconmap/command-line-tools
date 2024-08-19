# Reconmap command line tools

This monorepo contains the Reconmap [CLI](cli) and [agent](agent) command line tools, as well as the [shared library](shared-lib) written in Golang.

Look at each subdirectory to learn more about each tool including building and running instructions.

docker build -t quay.io/reconmap/agent:latest -f agent/Dockerfile .

## Build requirements

- Make
- Golang 1.22+
