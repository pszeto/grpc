FROM golang:alpine AS builder
WORKDIR /
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM golang:alpine
WORKDIR /
COPY --from=builder /main /main
COPY ca-cert.pem /cert/ca-cert.pem
COPY server-cert.pem /cert/server-cert.pem
COPY server-key.pem /cert/server-key.pem
EXPOSE 9000
ENTRYPOINT ["/main"]
