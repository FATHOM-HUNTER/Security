SHELL := /bin/bash

build:
	cd v2 && go mod download && pkger -o cmd/webserver && go build -o ../ ./cmd/webserver && cd ../

run: build
	chmod +x ./webserver
	# sudo setcap 'cap_net_bind_service=+ep' ./webserver
	./webserver