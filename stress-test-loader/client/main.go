package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	pb "stress-test-loader/proto"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
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

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	if os.Args[1] != "-s" && os.Args[1] != "-p" {
		// Start the stress test
		fmt.Println("Starting the stress test")
		loadTestConfig = os.Args[1]
		pbRequest = readStressTestConfig(loadTestConfig)
		pbRequest.TimeStamp = time.Now().UTC().Format(time.RFC3339Nano)
		callServer(&pbRequest, port, ipList, tlsCredentials, "start")
		fmt.Println("Stress test started")

		if len(os.Args) > 3 {
			maxTime, err := strconv.Atoi(os.Args[3])
			if err != nil {
				// Handle error if the string cannot be converted to an int
				fmt.Println("Error:", err)
				fmt.Println("Error: The third argument should be a positive integer which specifies the maximum time to run the stress tests in seconds")
				return
			}
			if maxTime < 0 {
				fmt.Println("Error: The third argument should be a positive integer which specifies the maximum time to run the stress tests in seconds")
				return
			}
			maxIteration := 10
			interval := float64(maxTime) / float64(maxIteration)
			// Poll the status of the stress test every interval seconds
			for i := 0; i < maxIteration; i++ {
				fmt.Println("Polling the status of the stress test")
				runningCount := callServer(&pbRequest, port, ipList, tlsCredentials, "poll")
				if runningCount == 0 {
					fmt.Println("All stress tests finished. Exiting")
					return
				} else {
					fmt.Println("Stress tests on", runningCount, "servers not finished. Sleeping for", interval, "seconds")
					time.Sleep(time.Duration(interval) * time.Second)
				}
			}

			fmt.Println("Stress tests are still running after", maxTime, "seconds. Stopping the tests")
			callServer(&pbRequest, port, ipList, tlsCredentials, "stop")

		}
	} else if os.Args[1] == "-p" {
		// Poll the status of the stress test
		fmt.Println("Getting the stress test status")
		runningCount := callServer(&pbRequest, port, ipList, tlsCredentials, "poll")
		fmt.Println("Running stress tests count:", runningCount)
	} else if os.Args[1] == "-s" {
		// Stop the stress test
		fmt.Println("stop stress test")
		callServer(&pbRequest, port, ipList, tlsCredentials, "stop")
	}
}

func callServer(pbRequest *pb.TestRequest, port string, ipList [][]PublicIP, tlsCredentials credentials.TransportCredentials, action string) (runningCount int) {

	runningCount = 0
	for _, s := range ipList {
		for _, s2 := range s {
			fmt.Println(s2.PublicIP)
			ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Second)
			conn, err := grpc.DialContext(ctx, s2.PublicIP+":"+port, grpc.WithTransportCredentials(tlsCredentials))
			if err != nil {
				log.Println("Dial failed!")
				return
			}
			defer conn.Close()
			c := pb.NewStressTestLoaderClient(conn)

			defer cancel()

			if action == "stop" {
				r, err := c.StopStressTest(ctx, pbRequest)
				if err != nil {
					log.Error("could not greet: %v", err)
				}
				log.Printf("%s", r.GetStatus())
			} else if action == "poll" {
				r, err := c.GetStressTestStatus(ctx, pbRequest)
				if err != nil {
					log.Error("could not greet: %v", err)
				}
				log.Printf("%s", r.GetStatus())
				if strings.Contains(r.GetStatus(), "running") {
					runningCount += 1
				}
			} else if action == "start" {
				r, err := c.StartStressTest(ctx, pbRequest)
				if err != nil {
					log.Error("could not greet: %v", err)
				}
				log.Printf("%s", r.GetStatus())
			}

		}
	}
	return runningCount
}
