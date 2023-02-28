.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
	# shadow ./... # this tool detects shadowing variables
.PHONY:vet

build: vet
	go generate ./ent
	go build
.PHONY:build