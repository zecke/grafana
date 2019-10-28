-include local/Makefile

.PHONY: all deps-go deps-js deps build-go build-server build-js build build-docker-dev build-docker-full lint-go gosec revive golangci-lint go-vet test-go test-js test run clean devenv devenv-down revive-alerting help

GO = GO111MODULE=on go
GO_FILES ?= ./pkg/...
SH_FILES ?= $(shell find ./scripts -name *.sh)

all: deps build

deps-go:
	$(GO) run build.go setup

deps-js: node_modules

deps: deps-js

node_modules: package.json yarn.lock
	@echo "install frontend dependencies"
	yarn install --pure-lockfile --no-progress

build-go:
	@echo "build go files"
	$(GO) run build.go build

build-server:
	@echo "build server"
	$(GO) run build.go build-server

build-js:
	@echo "build frontend"
	yarn run build

build: build-go build-js

scripts/go/bin/bra: scripts/go/go.mod
	@cd scripts/go; \
	$(GO) build -o ./bin/bra github.com/unknwon/bra

run: scripts/go/bin/bra ## Build and run web server on filesystem changes.
	@GO111MODULE=on scripts/go/bin/bra run

test-go:
	@echo "test backend"
	$(GO) test -v ./pkg/...

test-js:
	@echo "test frontend"
	yarn test

test: test-go test-js

scripts/go/bin/revive: scripts/go/go.mod
	@cd scripts/go; \
	$(GO) build -o ./bin/revive github.com/mgechev/revive

revive: scripts/go/bin/revive
	@echo "lint via revive"
	@scripts/go/bin/revive \
		-formatter stylish \
		-config ./scripts/go/configs/revive.toml \
		$(GO_FILES)

revive-alerting: scripts/go/bin/revive
	@echo "lint alerting via revive"
	@scripts/go/bin/revive \
		-formatter stylish \
		./pkg/services/alerting/...

scripts/go/bin/golangci-lint: scripts/go/go.mod
	@cd scripts/go; \
	$(GO) build -o ./bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

golangci-lint: scripts/go/bin/golangci-lint
	@echo "lint via golangci-lint"
	@scripts/go/bin/golangci-lint run \
		--config ./scripts/go/configs/.golangci.yml \
		$(GO_FILES)

scripts/go/bin/gosec: scripts/go/go.mod
	@cd scripts/go; \
	$(GO) build -o ./bin/gosec github.com/securego/gosec/cmd/gosec

gosec: scripts/go/bin/gosec
	@echo "lint via gosec"
	@scripts/go/bin/gosec -quiet \
		-exclude=G104,G107,G108,G201,G202,G204,G301,G304,G401,G402,G501 \
		-conf=./scripts/go/configs/gosec.json \
		$(GO_FILES)

go-vet:
	@echo "lint via go vet"
	@$(GO) vet $(GO_FILES)

lint-go: go-vet golangci-lint revive revive-alerting gosec ## Run all code checks for backend.

# with disabled SC1071 we are ignored some TCL,Expect `/usr/bin/env expect` scripts
shellcheck: $(SH_FILES) ## Run checks for shell scripts.
	@docker run --rm -v "$$PWD:/mnt" koalaman/shellcheck:stable \
	$(SH_FILES) -e SC1071 -e SC2162

build-docker-dev:
	@echo "build development container"
	@echo "\033[92mInfo:\033[0m the frontend code is expected to be built already."
	$(GO) run build.go -goos linux -pkg-arch amd64 ${OPT} build pkg-archive latest
	cp dist/grafana-latest.linux-x64.tar.gz packaging/docker
	cd packaging/docker && docker build --tag grafana/grafana:dev .

build-docker-full:
	@echo "build docker container"
	docker build --tag grafana/grafana:dev .

ifeq ($(sources),)
devenv:
	@printf 'You have to define sources for this command \nexample: make devenv sources=postgres,openldap\n'
else
devenv: devenv-down ## Start optional services, e.g. postgres, prometheus, and elasticsearch.
	$(eval targets := $(shell echo '$(sources)' | tr "," " "))

	@cd devenv; \
	./create_docker_compose.sh $(targets) || \
	(rm -rf {docker-compose.yaml,conf.tmp,.env}; exit 1)

	@cd devenv; \
	docker-compose up -d --build
endif

devenv-down:
	@cd devenv; \
	test -f docker-compose.yaml && \
	docker-compose down || exit 0;

clean:
	@echo "cleaning"
	rm -rf node_modules
	rm -rf public/build
