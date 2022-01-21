ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build


.PHONY: generate
generate: swagger
	mkdir -p $(shell pwd)/api/pipeline
	$(SWAGGER) generate client -f https://raw.githubusercontent.com/kubeflow/pipelines/master/backend/api/swagger/pipeline.swagger.json \
		--target $(shell pwd)/api/pipeline
	mkdir -p $(shell pwd)/api/pipeline_upload
	$(SWAGGER) generate client -f https://raw.githubusercontent.com/kubeflow/pipelines/master/backend/api/swagger/pipeline.upload.swagger.json \
		--target $(shell pwd)/api/pipeline_upload
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: generate fmt vet
	go test ./... -coverprofile cover.out

ifndef ignore-not-found
  ignore-not-found = false
endif

SWAGGER = $(shell pwd)/bin/swagger

.PHONY: swagger
swagger:
	$(call go-get-tool,$(SWAGGER),github.com/go-swagger/go-swagger/cmd/swagger@latest)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
