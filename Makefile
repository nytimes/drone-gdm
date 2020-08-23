#===============================================================================
#
# drone-gdm/makefile:
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


#=======================
# Makefile Parameters:
#-----------------------
# Makefile convenience vars:
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(shell dirname $(mkfile_path) )

#---------------------------
# output binary name is
# drone-gdm.<GOOS>.<GOARCH>
#
# Override with:
#   DRONE_GDM_BIN=whatever make
#
gdm_bin_default := drone-gdm
ifdef GOOS
	gdm_bin_default := $(gdm_bin_default).$(GOOS)
endif
ifdef GOARCH
	gdm_bin_default := $(gdm_bin_default).$(GOARCH)
endif

DRONE_GDM_BIN ?= $(gdm_bin_default)
#------------------------

# version name (tag, branch, etc):
DRONE_GDM_LABEL ?= LOCAL
# build revision:
DRONE_GDM_REVISION ?= $(shell date +%s)

# go build LDFLAGS:
DRONE_GDM_LDFLAGS ?= -X main.lbl=$(DRONE_GDM_LABEL) -X main.rev=$(DRONE_GDM_REVISION)


#=======================
# Targets:
#-----------------------
# PHONY targets:
# (See: https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html)
.PHONY: all vet test clean check bin

all: vet test bin

# Go vet
vet:
	go vet $(mkfile_dir)/cmd/drone-gdm

# Generic alias target for machine-specific binary:
bin: $(DRONE_GDM_BIN)

# Build drone-gdm for the current target:
$(DRONE_GDM_BIN): cmd/drone-gdm/*.go internal/plugin/*.go
	@echo 'Building "$(DRONE_GDM_BIN)":'
	@echo '- GOARCH="$(GOARCH)"'
	@echo '- GOOS="$(GOOS)"'
	@echo '- CGO_ENABLED="$(CGO_ENABLED)"'
	GOARCH=$(GOARCH) GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) \
	       go build -ldflags "$(DRONE_GDM_LDFLAGS)" \
	       -o $(DRONE_GDM_BIN) $(DRONE_GDM_BUILD_FLAGS) \
	       cmd/drone-gdm/main.go

# Execute tests:
test:
	go test -v ./...

# Alias for old-school make invocations:
check: test

# Clean up build output:
clean:
	rm -vf $(DRONE_GDM_BIN)

