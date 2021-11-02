SHELL := /bin/bash

build:
	cd frontend && yarn install && yarn build
	rm -rf static
	cp -r frontend/dist/ static
	go mod download
	rm -f v2/cmd/webserver/pkged.go
	pkger
	go build -o webserver

docker-build:
	docker build -t azurewebserver.azurecr.io/webserver .

docker-run: docker-build
	docker run -p 8081:8081 azurewebserver.azurecr.io/webserver

docker-push: 
	docker push azurewebserver.azurecr.io/webserver

run: build
	chmod +x ./webserver
	# sudo setcap 'cap_net_bind_service=+ep' ./webserver
	./webserver