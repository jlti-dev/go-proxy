FROM golang:1.17 as builder
WORKDIR /app
COPY go /app
RUN go mod download && CGO_ENABLED=1 GOOS=linux GOFLAGS=-mod=mod go build -a -installsuffix cgo -o main .

FROM debian:latest
WORKDIR /docker
RUN apt update && apt install -y iproute2
COPY --from=builder /app/main /docker/main

copy start.sh /docker/start.sh

EXPOSE 8080

CMD ["/bin/bash", "/docker/start.sh"]
