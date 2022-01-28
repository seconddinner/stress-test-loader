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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// tag the version of the code here
var Version = "notset"

var VrList pb.StressTestConfig

type server struct {
	pb.LoadTestLoaderServer
}

func copyStressTest(in *pb.TestRequest) (err error) {
	file, err := os.Create("/tmp/cli_s3.tgz")
	if err != nil {
		log.Error("Unable to open file %q, %v", in.S3Key, err)
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(in.S3),
			Key:    aws.String(in.S3Key),
		})
	if err != nil {
		log.Error("Unable to download in.S3Key %q, %v", in.S3Key, err)
	}

	log.Println("Downloaded", file.Name(), numBytes, "bytes")
	return
}

func startLoadTest(in *pb.TestRequest) (err error) {
	log.Print(in)

	err = copyStressTest(in)
	if err != nil {
		return
	}
	cmd := exec.Command(in.LoadtestExec)
	cmd.Env = os.Environ()
	for _, s := range in.EnvVariableList {
		cmd.Env = append(cmd.Env, s.EnvName+"="+s.EnvValue)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Error("cmd.Run: %s failed: %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 1 {
		log.Print("out:\n%s\nerr:\n%s\n", outStr, errStr)
	}

	log.Print(outStr)
	return
}

// not used right now, in the future, update config with this function dynamically
func (s *server) StartLoadTest(ctx context.Context, in *pb.TestRequest) (*pb.TestReply, error) {
	go startLoadTest(in)
	return &pb.TestReply{Status: "Hello again "}, nil
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
		if *VrList.DebugL == pb.StressTestConfig_DebugLevel {
			log.SetLevel(log.DebugLevel)
		}
	}

	// setup grpc server, this app is build for external config if needed in the future
	s := grpc.NewServer()
	pb.RegisterLoadTestLoaderServer(s, &server{})
	log.Printf("%s grpc server listening at %v.", Version, lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
