drone-gdm
=========

[![Build Status](https://travis-ci.org/NYTimes/drone-gdm.svg?branch=master)](https://travis-ci.org/NYTimes/drone-gdm)

A simple drone plugin which wraps [Google Deployment Manager](https://cloud.google.com/deployment-manager/docs/). For the latest
_stable release_, use the `v2-stable` tag. For more info on specific versions,
see [tags](#tags).

### Features
 * Set the desired `state` (absent, present, or latest) and the plugin determines whether to create, update, or delete.
 * Support for all GDM v1 types, composites, type-providers, and beta/alpha features

#### Simple Example

:information_source: see [examples](./doc/EXAMPLES.md) for more detiled examples.
```yaml
deploy:
  gdm:
    image: nytimes/drone-gdm:v2-stable
    # Provided JSON auth token (from drone secrets):
    token: >
      $$GOOGLE_JSON_CREDENTIALS
    project: my-gcp-project
    configurations:
    - name: my-deployment
      group: deployment
      state: latest
      description: A basic GDM deployment yaml file which creates some resources
      path: ./my-deployment.yaml

```

Usage
-----
The bulk of the input parameters are mapped directly to `gcloud` command options.
Documentation follows for the handful of parameters which are particular to `drone-gdm`.

### State and Action
The `state` can be one of `absent`, `present`, or `latest`.

| Plugin "state" | Object Exists? | Action      |
| -------------- | -------------- | ----------- |
| present        | no             | `create`    |
| present        | yes            | _no action_   |
| latest         | no             | `create`    |
| latest         | yes            | `update`    |
| absent         | no             | _no action_   |
| absent         | yes            | `delete`    |

The specific `action` selected by drone-gdm can be provided to your template
as a property, by specifying `passAction: true`. This will invoke your
configuration or template with `--properties=action:<action from table above>`.

### Variables
To circumvent data-type limitations imposed by the passing of properties via the
deployment manager `--properties` option, external configuration files (see the
[examples](./doc/EXAMPLES.md) for more info), are processed first as [golang templates](https://golang.org/pkg/text/template/) with the following top-level interfaces available for variable interpolation:
 - `.drone` - Drone environment variables provided by the CI system during plugin invocation
 - `.plugin` - Plugin parameters passed via environment during plugin invocation
 - `.context` - Any variables defined in the `vars` section of the plugin invocation
 - `.config` - Any variables defined in the `vars` section of the configuration definition
 - `.properties` - Variables defined in the `properties` section of the configuration definition
 - `.gdm` - A dictionary containing:
   - `name` - entity name for the configuration/template/composite
   - `status` - the entity status (e.g. DEPRECATED, EXPERIMENTAL, SUPPORTED)
   - `project` - the GCP project name
   - `action` - the gcloud "action" parameter (i.e. `create`, `update`, or `delete`)


Building
--------
This project uses [go dep](https://github.com/golang/dep) for depdenency management. Additionally, the
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

Tags
----

#### All Versions
* [specific release number](https://github.com/NYTimes/drone-gdm/releases) (see [dockerhub repo](https://hub.docker.com/r/nytimes/drone-gdm/tags/)).
* the `develop` tag to get the last thing that _built_

#### 2.x
Starting with version `2.0.0a` the tag scheme is prefixed with major version, e.g:
* the `v2-alpha` tag to get the latest 2.x _alpha_ release
* the `v2-beta` tag to get the latest 2.x _beta_ release
* the `v2-stable` tag to get the latest 2.x _stable_ release
<sub>This pattern will continue with subsequent major version releases; enabling you to pin your build to the latest stable version of any given backwards-compatible, major-level release</sub>

#### 1.x Series
* the `latest` tag to get the latest *v1.x* _stable_
* the `beta` tag to get the latest _beta_ release
* the `alpha` tag to get the latest `alpha` release
* the `develop` tag to get the last thing that _built_
<sub>(alpha, beta, and develop tags introduced as of `1.2.1a`)</sub>


### Resources
 - [drone-gdm on Travis-CI](https://travis-ci.org/NYTimes/drone-gdm)
 - [drone-gdm on dockerhub](https://hub.docker.com/r/nytimes/drone-gdm/)
