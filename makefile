DEFAULT_GOAL := build

build:
	sam build
.PHONY:build

# make start DEBUG=--debug
start:
	sam local start-api --env-vars env.json -p 9090 $(DEBUG)
.PHONY:start