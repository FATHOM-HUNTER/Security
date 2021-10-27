# Build Frontend
FROM node:16.12.0-alpine3.14 as frontenv

LABEL maintainer="Sahil Gupte <sahilsachingupte@gmail.com>"

RUN mkdir -p /build
WORKDIR /build
COPY frontend/ .

RUN npm rebuild node-sass
RUN yarn install && yarn build


# Build backend
FROM golang:alpine AS buildenv

LABEL maintainer="Sahil Gupte <sahilsachingupte@gmail.com>"

ARG VERSION
ARG GIT_COMMIT

ENV BIN=${BIN}
ENV VERSION=${VERSION}
ENV GIT_COMMIT=${GIT_COMMIT}
ENV CGO_ENABLED=0

# Create a location in the container for the source code.
RUN mkdir -p /app/v2

# Copy the module files first and then download the dependencies. If this
# doesn't change, we won't need to do this again in future builds.
COPY v2/go.* /app/v2/

WORKDIR /app/v2
RUN go mod download
RUN go mod verify

# Copy the source code into the container.
COPY v2/pkg pkg
COPY v2/cmd cmd
COPY --from=frontenv /build/dist static

RUN go install github.com/markbates/pkger/cmd/pkger
RUN pkger

RUN go build -o /go/bin/app
RUN ["chmod", "+x", "/go/bin/app"]

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

# Final build
FROM scratch
ENV DEBUG=FALSE
COPY --from=buildenv /go/bin/app /go/bin/app
COPY config.env config.env
COPY migrations migrations
COPY --from=buildenv /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 80
ENTRYPOINT ["/go/bin/app"]