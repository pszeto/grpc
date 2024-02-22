#/bin/sh
docker rmi -f pszeto/grpc-tls-client
docker build -t pszeto/grpc-tls-client -f Dockerfile . --platform linux/amd64 --no-cache
export CGO_ENABLED=0 && go build -o ./bin/grpc-client-arm64
export CGO_ENABLED=0 && export GOOS=linux && export GOARCH=amd64 && go build -o ./bin/grpc-client-amd64
docker push pszeto/grpc-tls-client