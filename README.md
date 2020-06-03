# drone-gdm

[![Build Status](https://travis-ci.org/NYTimes/drone-gdm.svg?branch=master)](https://travis-ci.org/NYTimes/drone-gdm)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE.md)
[![Docker Pulls](https://img.shields.io/docker/pulls/nytimes/drone-gdm)](https://hub.docker.com/r/nytimes/drone-gdm)

A simple drone plugin which wraps [Google Deployment Manager](https://cloud.google.com/deployment-manager/docs/), licensed under the [Apache 2.0 License](./LICENSE.md).

### Features
 * Set the desired `state` (absent, present, or latest) and the plugin determines whether to create, update, or delete.
 * Support for all GDM v1 types, composites, type-providers, and beta/alpha features

:information_source: See the [usage](./doc/USAGE.md) and [examples](./doc/EXAMPLES.md) for additional detail.

#### Example Usage

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


#### Versions/Compatibility

The latest stable release is [`v2-stable`](https://hub.docker.com/r/nytimes/drone-gdm/tags).
It is compatible with [drone](https://drone.io/) `0.5` - `0.8`.

> (:information_source: Drone `1.x` compatibility coming in `v3`).

## Resources
 - [Usage](./doc/USAGE.md) and [examples](./doc/EXAMPLES.md)
 - [Contributing guide](./CONTRIBUTING.md)
 - [Development guide - building, CI, and cutting releases](./doc/DEVELOPMENT.md)
 - [Releases](./CHANGELOG.md):
   - [github](https://github.com/NYTimes/drone-gdm/releases)
   - [dockerhub](https://hub.docker.com/r/nytimes/drone-gdm/tags/)

