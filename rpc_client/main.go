package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "../services"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "Earth"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	//name := defaultName
	/*if len(os.Args) > 1 {
		name = os.Args[1]
	}*/
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//r, err := c.GetBodyPosition(ctx, &pb.BodyName{Name: name})
	r, err := c.RequestSolarSystemStats(ctx, &pb.BodyName{Name: ""})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("%v", r)

	/*r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: &name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())*/
}
