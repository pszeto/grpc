package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/mrtkmynsndev/grpc-tls-go/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type greeterService struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *greeterService) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received name: %v", request.GetName())
	return &pb.HelloReply{Message: "Hello " + request.GetName()}, nil
}

func main() {
	address, present := os.LookupEnv("LISTEN_ADDRESS")
	if !present {
		log.Println("Listen address not defined.  Defaulting to 0.0.0.0")
		address = "0.0.0.0"
	}
	port, present := os.LookupEnv("PORT")
	if !present {
		log.Println("Port not defined.  Defaulting to 9000")
		port = "9000"
	}
	log.Println("Starting GRPC server : ", address+":"+port)
	log.Println("Version 0.1")
	// listen port
	lis, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		log.Fatalf("Error on listen err: %v", err)
	}
	tlsEnv, present := os.LookupEnv("TLS")
	secure := false
	if !present {
		log.Println("TLS not defined.  Defaulting to non-secure")
		secure = false
	} else {
		if strings.ToLower(tlsEnv) == "true" {
			secure = true
		}
	}

	var grpcServer *grpc.Server
	if secure {
		// read ca's cert, verify to client's certificate
		caPem, err := os.ReadFile("/cert/ca-cert.pem")
		if err != nil {
			log.Fatal(err)
		}

		// create cert pool and append ca's cert
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caPem) {
			log.Fatal(err)
		}

		// read server cert & key
		serverCert, err := tls.LoadX509KeyPair("/cert/server-cert.pem", "/cert/server-key.pem")
		if err != nil {
			log.Fatal(err)
		}

		// configuration of the certificate what we want to
		conf := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
		}

		//create tls certificate
		tlsCredentials := credentials.NewTLS(conf)

		// create grpc server
		grpcServer = grpc.NewServer(grpc.Creds(tlsCredentials))
	} else {
		grpcServer = grpc.NewServer()
	}

	// register service into grpc server
	pb.RegisterGreeterServiceServer(grpcServer, &greeterService{})
	log.Printf("Listening at %v", lis.Addr())

	if secure {
		log.Println("Mode: TLS")
	} else {
		log.Println("Mode: Un-secure")
	}
	// listen port
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve err: %v", err)
	}
}
