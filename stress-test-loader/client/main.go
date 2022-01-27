package main

import (
	"context"
	"log"
	"os"
	pb "stress-test-loader/proto"
	"time"

	"google.golang.org/grpc"
)

func main() {
	var port string
	port = "9005"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)

	conn, err := grpc.DialContext(ctx, "localhost:"+port, grpc.WithInsecure())
	if err != nil {
		log.Println("Dial failed!")
		return
	}

	defer conn.Close()
	c := pb.NewLoadTestLoaderClient(conn)

	defer cancel()

	r, err := c.StartLoadTest(ctx, &pb.TestRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetStatus())
}
