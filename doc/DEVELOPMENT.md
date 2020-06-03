# Drone-GDM: Development

## Building

This project uses [go dep](https://github.com/golang/dep) for dependency management. Additionally, the
[3rd party dependencies](./vendor) are _committed into the repo_. However, the usual commands
apply:
 - `dep ensure` - to make sure dependencies are up to date
 - `go vet` - ensure code sanity
 - `go build` - build for local platform
 - `go test -v ./...` - run test suite
 - `./util/docker-build.sh` - prepare an executable in AMD64/Linux format for docker packaging
 - `docker build -t drone-gdm:local ./` - package drone-gdm as a docker image

There is also a [makefile](./makefile) which provides some shortcuts - specifically:
 - `make` - build and test
 - `make drone-gdm` - _just build_ the executable, without testing
 - `make test` - just test the executable (though, it is built, if absent)
 - `make docker-bin` - leverage [docker-build.sh](./util/docker-build.sh) to build a dockerizable executable
 - `make clean` - clean any compiler generated output from the repo

## Travis CI

Drone-gdm uses [Travis CI](https://travis-ci.org/) to automate build, test, and deployment.

The [`.travis.yml`](../.travis.yml) configuration file is used to set up CI environment. It then
invokes:

 - [travis/script.sh](../travis/script.sh) to build and test the binary
 - [travis/after-success.sh](../travis/after-success.sh) to build, tag, and push the [docker image](https://hub.docker.com/r/nytimes/drone-gdm).


## Tags
The `develop` tag to get the last thing that _built_. Releases are tagged as follows:

#### 2.x
Starting with version `2.0.0a` the tag scheme is prefixed with major version, e.g:
* `v2-alpha`: latest 2.x _alpha_ release
* `v2-beta`: latest 2.x _beta_ release
* `v2-stable`: latest 2.x _stable_ release

#### 1.x Series
* `latest`: latest *v1.x* _stable_
* `beta`: latest _beta_ release
* `alpha`: latest `alpha` release
* `develop`: last thing that _built_

