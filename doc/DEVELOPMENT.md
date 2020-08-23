# Drone-GDM: Development

> :warning: Additional documentation for maintainers can be found  [here](./MAINTAINING.md).

### Contents

 - [Building and testing](#building-and-testing)
   - [Docker](#docker)
   - [Source Info](#sources)
 - [Travis CI/Docker Hub](#travis-ci)

----

## Building and Testing

This project uses [go mod](https://blog.golang.org/using-go-modules) for dependency management.

:information_source: There is a [makefile](../makefile) which provides some shortcuts - specifically:
 - `make` - build and test
 - `make bin` - _just build_ the executable, without testing
 - `make test` - just test the executable (it _is_ built, if absent)
 - `make clean` - clean any compiler generated output from the repo

### Docker

To build and package the docker image:

```bash
GOOS=linux GOARCH=amd64 make bin \
  && docker build -t drone-gdm:local .
```

### Sources
 - [main.go](../cmd/drone-gdm/main.go): Plugin entrypoint
 - [config.go](../internal/plugin/config.go): Encapsulates top-level plugin/GDM configuration
 - [context.go](../internal/plugin/context.go): `gcloud` execution context (path, global options, etc)
 - [run.go](../internal/plugin/run.go): Base functionality for executing `gcloud` and capturing output
 - [gdm.go](../internal/plugin/gdm.go): GDM command line arg formatting and execution (using [run.go](../internal/plugin/run.go))
 - [composite.go](../internal/plugin/composite.go): CLI options and validation particular to [composite types](https://cloud.google.com/deployment-manager/docs/fundamentals#composite_types)
 - [deployment.go](../internal/plugin/deployment.go): CLI options and validation particular to [deployments](https://cloud.google.com/deployment-manager/docs/fundamentals#deployment)
 - [typeprovider.go](../internal/plugin/typeprovider.go): CLI options options and validation particular to [type providers](https://cloud.google.com/deployment-manager/docs/fundamentals#basetypes)

#### Utilities
 - [drone.go](../internal/plugin/drone.go): Fetch drone-specific parameters from the environment at startup
 - [parse.go](../internal/plugin/parse.go): Drone parameter parser (from environment variables)
 - [yaml2json.go](../internal/plugin/yaml2json.go): Utility functions to ease some difficulties regarding parsing YAML/JSON from strings

#### Tests
 - [parse_test.go](../internal/plugin/parse_test.go)
 - [run_test.go](../internal/plugin/run_test.go)

----

## Travis CI

Drone-gdm uses [Travis CI](https://travis-ci.org/) to automate build, test, and deployment.

The [`.travis.yml`](../.travis.yml) configuration file is used to set up CI environment. It then
invokes:

 - [travis/script.sh](../travis/script.sh) to build and test the binary
 - [travis/after-success.sh](../travis/after-success.sh) to build, tag, and push the [docker image](https://hub.docker.com/r/nytimes/drone-gdm).

