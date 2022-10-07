FROM golang:1.19 AS build
WORKDIR /src
RUN apt update 
RUN apt install ca-certificates gcc libzmq3-dev -y
COPY . .
RUN go mod vendor
RUN  GOOS=linux go build -ldflags="-w -s" -o matrix .
FROM ubuntu
RUN apt update && apt install libzmq3-dev -y && apt-get clean autoclean \
    && apt-get autoremove --yes \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/

COPY --from=build /src/matrix /opt/matrix
ENTRYPOINT ["/opt/matrix"]