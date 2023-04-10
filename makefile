.PHONY: build

build:
	sam build

start:
	# make start DEBUG=--debug
	sam local start-api --env-vars env.json -p 9091 $(DEBUG)