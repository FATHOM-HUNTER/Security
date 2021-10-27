SHELL := /bin/bash

build:
	cd frontend && yarn install && yarn build
	rm -rf static
	cp -r frontend/dist/ static
	go mod download
	rm -f v2/cmd/webserver/pkged.go
	pkger
	go build -o webserver

run: build
	chmod +x ./webserver
	# sudo setcap 'cap_net_bind_service=+ep' ./webserver
	./webserver