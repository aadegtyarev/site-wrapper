.PHONY: check

check:
	docker run --rm --user $$(id -u):$$(id -g) --env GOCACHE=/tmp/gocache --env GOPATH=/tmp/gopath --volume "$$(pwd):/src" --workdir /src golang:1.25 sh -c 'test -z "$$(gofmt -l internal)" && go vet ./... && go test ./...'
