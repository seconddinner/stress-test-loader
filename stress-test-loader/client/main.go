package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	pb "stress-test-loader/proto"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type PublicIP struct {
	PublicIP string `json:"public_ip"`
}

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

	var loadTestConfig string
	var pbRequest pb.TestRequest
	port := "9005"
	var ipList [][]PublicIP

	if len(os.Args) > 2 {

		jsonFile, err := os.Open(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &ipList)
	} else {
		var localhost PublicIP
		localhost.PublicIP = "localhost"
		var l1 []PublicIP
		l1 = append(l1, localhost)

		ipList = append(ipList, l1)
	}
	fmt.Println(ipList)

	if os.Args[1] != "-s" {
		loadTestConfig = os.Args[1]
		pbRequest = readStressTestConfig(loadTestConfig)

		for _, s := range ipList {
			for _, s2 := range s {
				fmt.Println(s2.PublicIP)
				ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
				conn, err := grpc.DialContext(ctx, s2.PublicIP+":"+port, grpc.WithInsecure())
				if err != nil {
					log.Println("Dial failed!")
					return
				}
				defer conn.Close()
				c := pb.NewStressTestLoaderClient(conn)

				defer cancel()

				r, err := c.StartStressTest(ctx, &pbRequest)
				if err != nil {
					log.Error("could not greet: %v", err)
				}
				log.Printf("Greeting: %s", r.GetStatus())

			}
		}
	} else {
		for _, s := range ipList {
			for _, s2 := range s {
				fmt.Println(s2.PublicIP)
				ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
				conn, err := grpc.DialContext(ctx, s2.PublicIP+":"+port, grpc.WithInsecure())
				if err != nil {
					log.Println("Dial failed!")
					return
				}
				defer conn.Close()
				c := pb.NewStressTestLoaderClient(conn)

				defer cancel()

				r, err := c.StopStressTest(ctx, &pbRequest)
				if err != nil {
					log.Error("could not greet: %v", err)
				}
				log.Printf("Greeting: %s", r.GetStatus())

			}

			fmt.Println("stop stress test")
		}
	}
}
