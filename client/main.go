package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/mrtkmynsndev/grpc-tls-go/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tlsFlag       = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	serverAddress = flag.String("server", "127.0.0.1", "The server address")
	port          = flag.Int("port", 9000, "The server port")
	name          = flag.String("name", "Peter", "name to say hello to")
)

func main() {
	flag.Parse()
	server := *serverAddress + ":" + strconv.Itoa(*port)
	var opts []grpc.DialOption
	if *tlsFlag {
		// read ca's cert
		caCert, err := os.ReadFile("/cert/ca-cert.pem")
		if err != nil {
			log.Fatal("Failed to read ca: /cert/ca-cert.pem")
			log.Fatal(caCert)
		}

		// create cert pool and append ca's cert
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(caCert); !ok {
			log.Fatal("Failed to create cert pool from ca: /cert/ca-cert.pem")
			log.Fatal(err)
		}

		//read client cert
		clientCert, err := tls.LoadX509KeyPair("/cert/client-cert.pem", "/cert/client-key.pem")
		if err != nil {
			log.Fatal("Failed to load client cert: /cert/client-cert.pem & /cert/client-key.pem")
			log.Fatal(err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      certPool,
		}

		tlsCredential := credentials.NewTLS(config)

		opts = append(opts, grpc.WithTransportCredentials(tlsCredential))
		log.Println("Creating tls connection")
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		log.Println("Creating in-secure connection")
	}
	opts = append(opts, grpc.WithAuthority("grpc-tls-server"))
	log.Println("Calling with :authority 'grpc-tls-server'")
	conn, err := grpc.Dial(
		server,
		opts...)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGreeterServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Greeting: %s", resp.GetMessage())
}
