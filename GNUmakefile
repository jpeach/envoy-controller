RM_F := rm -rf
GO := go
GIT := git
DOCKER := docker

export GO111MODULE=on

BIN := envoy-controller

REPO := github.com/jpeach/envoy-controller
SHA := $(shell git rev-parse --short=8 HEAD)
VERSION := $(shell git describe --exact-match 2>/dev/null || basename $$(git describe --all --long 2>/dev/null))
BUILDDATE := $(shell TZ=GMT date '+%Y-%m-%dT%R:%S%z')

GO_BUILD_LDFLAGS := \
	-s \
	-w \
	-X $(REPO)/pkg/version.Progname=$(BIN) \
	-X $(REPO)/pkg/version.Version=$(VERSION) \
	-X $(REPO)/pkg/version.Sha=$(SHA) \
	-X $(REPO)/pkg/version.BuildDate=$(BUILDDATE)
 
# Image URL to use all building/pushing image targets
IMG ?= $(BIN):$(VERSION)

.PHONY: help
help:
	@echo "$(BIN)"
	@echo
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: build
build: ## Build
build:
	@$(GO) build -mod=readonly -ldflags "$(GO_BUILD_LDFLAGS)" -o $(BIN) .

.PHONY: generate
generate: ## Generate build files
generate:
	$(GO) mod vendor
	for prog in ./hack/generate-* ; do \
		$$prog ; \
	done


.PHONY: check
check: ## Run tests
check:
	$(GO) test ./... -coverprofile cover.out


# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -


# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

.PHONY: lint
lint: ## Run linters
lint:
	@if command -v golangci-lint > /dev/null 2>&1 ; then \
		golangci-lint run --exclude-use-default=false ; \
	else \
		$(DOCKER) run \
			--rm \
			--volume $$(pwd):/app \
			--workdir /app \
			--env GO111MODULE \
			golangci/golangci-lint:v1.23.7 \
			golangci-lint run --exclude-use-default=false ; \
	fi

# Build the docker image
docker-build:
	$(DOCKER) build . -t ${IMG}

# Push the docker image
docker-push:
	$(DOCKER) push ${IMG}

.PHONY: clean
clean: ## Remove build artifacts
	$(GO) clean ./...
	$(RM_F) $(BIN) cover.out vendor/

