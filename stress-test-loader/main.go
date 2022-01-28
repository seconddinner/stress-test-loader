package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

var StressTestLoaderConfig pb.StressTestConfig

type server struct {
	pb.StressTestLoaderServer
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, 0755)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func copyStressTest(in *pb.TestRequest) (err error) {
	file, err := os.Create(StressTestLoaderConfig.WorkingFolder + "/" + in.S3Key)
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

	err = os.RemoveAll(StressTestLoaderConfig.WorkingFolder + "/bin")
	if err != nil {
		log.Error(err)
	}

	err = ensureDir(StressTestLoaderConfig.WorkingFolder + "/bin")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print(StressTestLoaderConfig.WorkingFolder + "/bin exist!")
	}
	cmd := exec.Command("tar", "xzf", file.Name(), "-C", StressTestLoaderConfig.WorkingFolder+"/bin")

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

func startStressTest(in *pb.TestRequest) (err error) {
	log.Print(in)

	err = copyStressTest(in)
	if err != nil {
		return
	}
	cmd := exec.Command(StressTestLoaderConfig.WorkingFolder + "/bin" + "/" + in.LoadtestExec)
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
func (s *server) StartStressTest(ctx context.Context, in *pb.TestRequest) (*pb.TestReply, error) {
	go startStressTest(in)
	return &pb.TestReply{Status: in.S3Key + " stress-test started"}, nil
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
	json.Unmarshal(byteValue, &StressTestLoaderConfig)
	jsonBytes, err := json.MarshalIndent(StressTestLoaderConfig, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(string(jsonBytes))
	err = ensureDir(StressTestLoaderConfig.WorkingFolder)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print(StressTestLoaderConfig.WorkingFolder + "config exist!")
	}

}

func main() {
	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(StressTestLoaderConfig.ListenPort), 10))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if StressTestLoaderConfig.DebugL != nil {
		// to do add more log levels
		if *StressTestLoaderConfig.DebugL == pb.StressTestConfig_DebugLevel {
			log.SetLevel(log.DebugLevel)
		}
	}

	// setup grpc server, this app is build for external config if needed in the future
	s := grpc.NewServer()
	pb.RegisterStressTestLoaderServer(s, &server{})
	log.Printf("%s grpc server listening at %v.", Version, lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
