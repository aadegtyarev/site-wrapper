.PHONY: check go-check extension-check desktop-package android-check

check: go-check extension-check

go-check:
	docker run --rm --user $$(id -u):$$(id -g) --env GOCACHE=/tmp/gocache --env GOPATH=/tmp/gopath --volume "$$(pwd):/src" --workdir /src golang:1.25 sh -c 'test -z "$$(gofmt -l internal cmd)" && go vet ./... && go test ./...'

extension-check:
	docker run --rm --user $$(id -u):$$(id -g) --env npm_config_cache=/tmp/npm-cache --volume "$$(pwd):/src" --workdir /src node:22 sh -c 'node --test extension/*.test.mjs && npm exec --yes --package=web-ext@10.8.0 -- web-ext lint --source-dir extension --self-hosted --warnings-as-errors'

desktop-package:
	mkdir -p dist
	docker run --rm --user $$(id -u):$$(id -g) --env GOCACHE=/tmp/gocache --env GOPATH=/tmp/gopath --volume "$$(pwd):/src" --workdir /src golang:1.25 go build -o dist/site-wrapper ./cmd/site-wrapper
	docker run --rm --user $$(id -u):$$(id -g) --volume "$$(pwd):/src" --workdir /src python:3.13-alpine python scripts/package-extension.py dist/site-wrapper@aadegtyarev.github.io.xpi

android-check:
	bash scripts/android-docker.sh
