SHELL := /bin/bash

build:
	cd v2 && go mod download && pkger -o cmd/webserver && go build -o ../ ./cmd/webserver && cd ../

run: build
	chmod +x ./webserver ./setupenv.sh
	sudo setcap 'cap_net_bind_service=+ep' ./webserver
	source setupenv.sh && ./webserver