.PHONY: build

build:
	sam build

start:
	sam local start-api --env-vars env.json -p 9090

start --debug:
	sam local start-api --env-vars env.json -p 9090 --debug