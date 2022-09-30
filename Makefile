SHELL := /bin/bash

build:
	cd frontend && yarn install && yarn build
	rm -rf static
	cp -r frontend/dist/ static
	go mod download
	rm -f v2/cmd/webserver/pkged.go
	pkger
	go build -o webserver

docker-build: build
	docker build -t specialization1.azurecr.io/webserver .

docker-run: docker-build
	docker run -p 8081:80 specialization1.azurecr.io/webserver

docker-push: docker-build
	docker push specialization1.azurecr.io/webserver

run:
	go run main.go