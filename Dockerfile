FROM golang:1.19 AS build
WORKDIR /src
RUN apt update && apt install ca-certificates && update-ca-certificates
RUN apt install gcc libzmq3-dev -y
COPY . .
RUN go mod vendor
RUN  GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o matrix .
ENTRYPOINT ["/src/matrix"]