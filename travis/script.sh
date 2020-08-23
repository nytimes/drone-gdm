#!/bin/bash
#===============================================================================
#
# drone-gdm/travis/script.sh:
#   - Run unit tests
#   - Build the binary for a target compatible with our docker image.
#
#-------------------------------------------------------------------------------
DRONE_GDM_LABEL="${TRAVIS_BRANCH}" \
    DRONE_GDM_BUILD="${TRAVIS_BRANCH}" \
    DRONE_GDM_REVISION="${TRAVIS_COMMIT}" \
    DRONE_GDM_BUILD_FLAGS="-a -tags netgo" \
    GOOS="linux" \
    GOARCH="amd64" \
    CGO_ENABLED=0 \
    make clean vet test bin

# EOF
