#!/bin/bash
GDM_LDFLAGSFLAGS=(
    "-s -w"
    "-X main.lbl=$TRAVIS_BRANCH"
    "-X main.build=$TRAVIS_BUILD_NUMBER"
    "-X main.rev=$TRAVIS_COMMIT"
    )
go vet
go test github.com/nytimes/drone-gdm/plugin
go build -ldflags "${GDM_LDFLAGSFLAGS[@]}" -a -tags netgo

# EOF
