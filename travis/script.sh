#!/bin/bash
GDM_LDFLAGS=""
GDM_LDFLAGS+="-s -w "
GDM_LDFLAGS+="-X main.lbl=$TRAVIS_BRANCH "
GDM_LD_FLAGS+="-X main.build=$TRAVIS_BUILD_NUMBER "
GDM_LD_FLAGS+=" -X main.rev=$TRAVIS_COMMIT"

go vet
go test ./...
go build -ldflags "${GDM_LDFLAGS}" -a -tags netgo

# EOF
