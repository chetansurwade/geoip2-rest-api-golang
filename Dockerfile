FROM golang:1.11-alpine AS build-env
WORKDIR /go/src/app
COPY . .

RUN apk update && apk add git \
    && go get \
    && go build -o goapp

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build-env /go/src/app/goapp .
COPY . .
EXPOSE 8080
ENTRYPOINT [ "./goapp" ]