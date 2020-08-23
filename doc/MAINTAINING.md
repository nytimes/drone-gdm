# Drone-GDM: Maintaining

This document provides instructions for drone-gdm maintainers. For general
development instructions, see [DEVELOPMENT.md](./DEVELOPMENT.md).

## Releases

Releases are cut using git tags. The tags use a modified subset
of [semantic versioning](https://semver.org/), as indicated below.

**TL;DR**: Releases are in the form `M.m.p[a|b]`, i.e. the semantic major, minor,
and patch numbers and - _optionally_ - a prerelease suffix (which must be either
`a` or `b`).

> :grimacing: Semver compliance is on the TODO list.

### Cutting a New Release

To cut a new release:
1. [Build and test locally](./DEVELOPMENT.md)
1. `git push origin main`
1. [Verify that the build passed](https://travis-ci.org/github/NYTimes/drone-gdm)
1. [Verify that the `develop` docker tag](https://hub.docker.com/r/nytimes/drone-gdm/tags) was updated
1. Update the [CHANGELOG](../CHANGELOG.md) with relevant info
1. `git tag -a <version>`
1. `git push origin main <version>`
1. Verify that the [release](../README.md#docker-tags) docker tag was updated

### Tag Format Spec

```ABNF
<valid tag> ::= <version core>
              | <version core><pre-release>

<version core> ::= <major> "." <minor> "." <patch>

<major> ::= <digits>

<minor> ::= <digits>

<patch> ::= <digits>

<digits> ::= <digit>
           | <digit> <digits>

<digit> ::= 0-9

<pre-release> ::= "a" | "b"
```

