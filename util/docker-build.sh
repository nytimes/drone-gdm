#!/usr/bin/env bash -e
#===============================================================================
#
# drone-gdm/docker-build.sh:
#     build drone-gdm with the appropriate configurations/target architecture
#     for packaging as a docker image.
#
# Copyright (c) 2017 The New York Times Company
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this library except in compliance with the License.
# You may obtain a copy of the License at:
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#-------------------------------------------------------------------------------
__thispath="${0}"
__thisdir="${__thispath%/*}"
__thisname="${__thispath##*/}"

DRONE_GDM_DIR="${DRONE_GDM_DIR:-"${__thisdir}/.."}"
cd "${DRONE_GDM_DIR}"

if [ -z "${GDM_LDFLAGS}" ]; then
    GDM_LDFLAGS=""
    GDM_LDFLAGS+=" -s -w"
    GDM_LDFLAGS+=" -X main.lbl=${TRAVIS_BRANCH}"
    GDM_LDFLAGS+=" -X main.build=${TRAVIS_BUILD_NUMBER}"
    GDM_LDFLAGS+=" -X main.rev=${TRAVIS_COMMIT}"
fi

export GOARCH="${DRONE_GDM_ARCH:-"amd64"}"
export GOOS="${DRONE_GDM_OS:-"linux"}"

# Print the value of a variable with name label.
# $1     - Variable name as a string
function print_var() {
    local var_name
    var_name="$1"
    printf "\e[00;34m%s\e[00m=\"\e[00;33m%s\e[00m\"\n" "${var_name}" "${!var_name}"
}

# Run a command, logging the path and arguments
function run_cmd() {
    local cmd
    cmd="$1" ; shift
    echo -e "\e[00;02mEXEC: \e[00;32m${cmd} $@\e[00m" >&2
    ${DRY} ${cmd} "$@"
}

function build_drone_gdm() {
    run_cmd go vet
    run_cmd go build -v
}

function main() {
    print_var "GDM_LDFLAGS"
    print_var "GOARCH"
    print_var "GOOS"

    build_drone_gdm
}

main "$@"

# EOF

