package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	pb "stress-test-loader/proto"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func readStressTestConfig(f string) (pbRequest pb.TestRequest) {
	// open config.json defined by protobuf
	jsonFile, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully Opened users.json")
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, &pbRequest)
	jsonBytes, err := json.MarshalIndent(pbRequest, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(string(jsonBytes))
	return
}

func main() {
	var host string
	var loadTestConfig string
	var pbRequest pb.TestRequest
	host = "localhost:9005"
	if len(os.Args) > 2 {
		host = os.Args[2]
	}
	loadTestConfig = os.Args[1]
	pbRequest = readStressTestConfig(loadTestConfig)
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)

	conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure())
	if err != nil {
		log.Println("Dial failed!")
		return
	}

	defer conn.Close()
	c := pb.NewStressTestLoaderClient(conn)

	defer cancel()

	r, err := c.StartStressTest(ctx, &pbRequest)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetStatus())
}
