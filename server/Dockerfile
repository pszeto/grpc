FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o grpc-server

FROM --platform=linux/amd64 alpine:3.19
WORKDIR /
COPY --from=builder /grpc-server /grpc-server
COPY ca-cert.pem /cert/ca-cert.pem
COPY server-cert.pem /cert/server-cert.pem
COPY server-key.pem /cert/server-key.pem
EXPOSE 9000
ENTRYPOINT ["/grpc-server"]
