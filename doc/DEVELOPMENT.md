# Drone-GDM: Development

> :warning: Additional documentation for maintainers can be found  [here](./MAINTAINING.md).

### Contents

 - [Building and testing](#building-and-testing)
   - [Dependency Management](#dependency-management)
   - [Docker](#docker)
   - [Source Info](#sources)
 - [Travis CI/Docker Hub](#travis-ci)

----

## Building and Testing

:information_source: There is a [makefile](../makefile) which provides some shortcuts - specifically:
 - `make` - build and test
 - `make drone-gdm` - _just build_ the executable, without testing
 - `make test` - just test the executable (it _is_ built, if absent)
 - `make docker-bin` - leverage [docker-build.sh](../util/docker-build.sh) to build a dockerizable executable
 - `make clean` - clean any compiler generated output from the repo

### Dependency Management
:information_source: This project uses [go dep](https://github.com/golang/dep) for dependency
management - [3rd party dependencies](../vendor) are _committed into the repo_.

> :construction: We are in the process of migrating to [go mod](https://blog.golang.org/using-go-modules).

The usual commands apply:
 - `dep ensure` - to make sure dependencies are up to date
 - `go vet` - ensure code sanity
 - `go build` - build for local platform
 - `go test -v ./...` - run test suite

### Docker
 - [`docker-build.sh`](../util/docker-builder.sh) - prepare an executable in AMD64/Linux format for docker packaging
 - `docker build -t drone-gdm:local ./` - package drone-gdm as a docker image

### Sources
 - [main.go](../main.go): Plugin entrypoint
 - [config.go](../config.go): Encapsulates top-level plugin/GDM configuration
 - [context.go](../context.go): `gcloud` execution context (path, global options, etc)
 - [run.go](../run.go): Base functionality for executing `gcloud` and capturing output
 - [gdm.go](../gdm.go): GDM command line arg formatting and execution (using [run.go](../run.go))
 - [composite.go](../composite.go): CLI options and validation particular to [composite types](https://cloud.google.com/deployment-manager/docs/fundamentals#composite_types)
 - [deployment.go](../deployment.go): CLI options and validation particular to [deployments](https://cloud.google.com/deployment-manager/docs/fundamentals#deployment)
 - [typeprovider.go](../typeprovider.go): CLI options options and validation particular to [type providers](https://cloud.google.com/deployment-manager/docs/fundamentals#basetypes)

#### Utilities
 - [drone.go](../drone.go): Fetch drone-specific parameters from the environment at startup
 - [parse.go](../parse.go): Drone parameter parser (from environment variables)
 - [yaml2json.go](../yaml2json.go): Utility functions to ease some difficulties regarding parsing YAML/JSON from strings

#### Tests
 - [parse_test.go](../parse_test.go)
 - [run_test.go](../run_test.go)

----

## Travis CI

Drone-gdm uses [Travis CI](https://travis-ci.org/) to automate build, test, and deployment.

The [`.travis.yml`](../.travis.yml) configuration file is used to set up CI environment. It then
invokes:

 - [travis/script.sh](../travis/script.sh) to build and test the binary
 - [travis/after-success.sh](../travis/after-success.sh) to build, tag, and push the [docker image](https://hub.docker.com/r/nytimes/drone-gdm).

