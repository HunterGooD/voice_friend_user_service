FROM golang:1.22-alpine as builder
WORKDIR /build
COPY go.mod . 
RUN go mod download && \
    go mod tidy
COPY . .
RUN go build -o /main cmd/grpc_server/main.go

FROM alpine:3 as runtime
COPY --from=builder main /bin/main
COPY --from=builder /build/config/config.yaml /etc/service/config.yaml
ENV LOG_LEVEL=debug
ENV CONFIG_PATH=/etc/service/config.yaml
ENTRYPOINT ["/bin/main"]
