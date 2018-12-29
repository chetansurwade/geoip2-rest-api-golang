FROM golang:1.11-alpine
WORKDIR /go/src/app
COPY . .

RUN apk update && apk add git \
    && go get \
    && go build main.go
EXPOSE 8080
CMD ["/go/src/app/main"]