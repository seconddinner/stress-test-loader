package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"

	pb "stress-test-loader/proto"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// tag the version of the code here
var Version = "notset"

var VrList pb.MonitorServerConfig

type server struct {
	pb.VerifierServer
}

func shutdownLoaderNode(shutdownFlag chan bool) {
	log.Printf("shutdown started")
	c := <-shutdownFlag
	if c == true {
		log.Printf("shutdown signal is true")
	}
	// call a shell command to log issued  host into an remote db.
	cmd := exec.Command("reportbadserver.sh")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("cmd.Run: %s failed: %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 1 {
		log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	}

	log.Printf(outStr)

	// If config.json enable shutdown flag, we will shutdown the
	if *VrList.ShutdownFlag == true {
		cmd = exec.Command("systemctl", "stop", "something")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		//  comment out the following part to check if this is causing stabliblites issue

		err = cmd.Run()
		if err != nil {
			log.Printf("cmd.Run: %s failed: %s\n", err)
		}
		outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())
		if len(errStr) > 1 {
			log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
		}

		log.Printf(outStr)

	}
	log.Printf("finished stopping")
}

// not used right now, in the future, update config with this function dynamically
func (s *server) Verify(ctx context.Context, in *pb.MonitorServerConfig) (*pb.VerifyReply, error) {
	return &pb.VerifyReply{Status: "Hello again "}, nil
}

// func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
// 	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
// }

func init() {
	// open config.json defined by protobuf
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully Opened users.json")
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, &VrList)
	jsonBytes, err := json.MarshalIndent(VrList, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(string(jsonBytes))
}

func main() {
	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(VrList.ListenPort), 10))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if VrList.DebugL != nil {
		// to do add more log levels
		if *VrList.DebugL == pb.MonitorServerConfig_DebugLevel {
			log.SetLevel(log.DebugLevel)
		}
	}

	// setup grpc server, this app is build for external config if needed in the future
	s := grpc.NewServer()
	pb.RegisterVerifierServer(s, &server{})
	log.Printf("%s grpc server listening at %v.", Version, lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
