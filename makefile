.DEFAULT_GOAL := build

go: # run the app with nodemon for hot reload
	nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run ./main.go
.PHONY:go

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