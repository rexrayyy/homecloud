FROM golang:1.21 AS build
WORKDIR /src
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /homecloud ./cmd
FROM ubuntu:22.04
COPY --from=build /homecloud /usr/local/bin/homecloud
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/homecloud", "--mode=server", "--addr=:8080"]