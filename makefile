include .env
GIT_REF_TAG := $(shell git describe --always --tag)
GIT_REF_COMMIT := $(shell git rev-parse --short=7 HEAD)

VERSION ?= $(GIT_REF_TAG)

.PHONY: lint
lint: ## lint all source codes
	@golangci-lint run ./...

.PHONY: build-app
build-app:
	@echo "+ $@"
	@echo "version=${VERSION}"
	go build -v \
        -ldflags="-w -s -X github.com/thoriqulumar/halo-suster/version.Version=${VERSION}" \
		-o bin/thorumr_halo-suster \
		cmd/main.go

.PHONY: run-app
run-app: 
	$(MAKE) build-app && ./bin/thorumr_halo-suster

.PHONY: migrate-up
migrate-up:
	migrate -path database/migration/ -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path database/migration/ -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down