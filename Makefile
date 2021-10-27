SHELL := /bin/bash

build:
	cd frontend && yarn install && yarn build
	rm -rf v2/static
	cp -r frontend/dist/ v2/static
	cd v2 && go mod download
	rm -f v2/cmd/webserver/pkged.go
	cd v2 && pkger
	cd v2 && go build -o ../webserver

run: build
	chmod +x ./webserver
	# sudo setcap 'cap_net_bind_service=+ep' ./webserver
	./webserver