FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o grpc-client

FROM --platform=linux/amd64 alpine:3.19
WORKDIR /
COPY --from=builder /grpc-client /grpc-client
COPY ca-cert.pem /cert/ca-cert.pem
COPY client-cert.pem /cert/client-cert.pem
COPY client-key.pem /cert/client-key.pem
EXPOSE 9000
CMD ["/bin/sleep infinity"]