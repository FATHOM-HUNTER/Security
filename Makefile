build:
    cd v2 && go build -o ../ ./cmd/webserver && cd ../

run: build
    chmod +x ./webserver ./setupenv.sh
    ./setupenv.sh && ./webserver