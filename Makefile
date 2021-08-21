run:
        build
        chmod +x ./webserver
        . setupenv.sh && ./webserver

build:
        cd v2 && go build ./cmd/webserver && cd ../