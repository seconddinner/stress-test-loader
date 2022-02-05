// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: proto/stresstestloader.proto

package stresstestloader

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StressTestConfig_DebugLevelEnum int32

const (
	StressTestConfig_PanicLevel StressTestConfig_DebugLevelEnum = 0
	StressTestConfig_FatalLevel StressTestConfig_DebugLevelEnum = 1
	StressTestConfig_ErrorLevel StressTestConfig_DebugLevelEnum = 2
	StressTestConfig_WarnLevel  StressTestConfig_DebugLevelEnum = 3
	StressTestConfig_InfoLevel  StressTestConfig_DebugLevelEnum = 4
	StressTestConfig_DebugLevel StressTestConfig_DebugLevelEnum = 5
	StressTestConfig_TraceLevel StressTestConfig_DebugLevelEnum = 6
)

// Enum value maps for StressTestConfig_DebugLevelEnum.
var (
	StressTestConfig_DebugLevelEnum_name = map[int32]string{
		0: "PanicLevel",
		1: "FatalLevel",
		2: "ErrorLevel",
		3: "WarnLevel",
		4: "InfoLevel",
		5: "DebugLevel",
		6: "TraceLevel",
	}
	StressTestConfig_DebugLevelEnum_value = map[string]int32{
		"PanicLevel": 0,
		"FatalLevel": 1,
		"ErrorLevel": 2,
		"WarnLevel":  3,
		"InfoLevel":  4,
		"DebugLevel": 5,
		"TraceLevel": 6,
	}
)

func (x StressTestConfig_DebugLevelEnum) Enum() *StressTestConfig_DebugLevelEnum {
	p := new(StressTestConfig_DebugLevelEnum)
	*p = x
	return p
}

func (x StressTestConfig_DebugLevelEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StressTestConfig_DebugLevelEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_stresstestloader_proto_enumTypes[0].Descriptor()
}

func (StressTestConfig_DebugLevelEnum) Type() protoreflect.EnumType {
	return &file_proto_stresstestloader_proto_enumTypes[0]
}

func (x StressTestConfig_DebugLevelEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StressTestConfig_DebugLevelEnum.Descriptor instead.
func (StressTestConfig_DebugLevelEnum) EnumDescriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{2, 0}
}

type EnvVariable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnvName  string `protobuf:"bytes,1,opt,name=EnvName,proto3" json:"EnvName,omitempty"`
	EnvValue string `protobuf:"bytes,2,opt,name=EnvValue,proto3" json:"EnvValue,omitempty"`
}

func (x *EnvVariable) Reset() {
	*x = EnvVariable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_stresstestloader_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvVariable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvVariable) ProtoMessage() {}

func (x *EnvVariable) ProtoReflect() protoreflect.Message {
	mi := &file_proto_stresstestloader_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvVariable.ProtoReflect.Descriptor instead.
func (*EnvVariable) Descriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{0}
}

func (x *EnvVariable) GetEnvName() string {
	if x != nil {
		return x.EnvName
	}
	return ""
}

func (x *EnvVariable) GetEnvValue() string {
	if x != nil {
		return x.EnvValue
	}
	return ""
}

// The request to StartStressTest
type TestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LoadtestExec    string         `protobuf:"bytes,1,opt,name=loadtestExec,proto3" json:"loadtestExec,omitempty"`
	S3              string         `protobuf:"bytes,2,opt,name=s3,proto3" json:"s3,omitempty"`
	EnvVariableList []*EnvVariable `protobuf:"bytes,3,rep,name=envVariableList,proto3" json:"envVariableList,omitempty"`
	// will add   bool runUntilStoped ;
	NumberOfGames        int32  `protobuf:"varint,4,opt,name=numberOfGames,proto3" json:"numberOfGames,omitempty"`
	S3Key                string `protobuf:"bytes,5,opt,name=s3key,proto3" json:"s3key,omitempty"`
	StepTestStart        int32  `protobuf:"varint,6,opt,name=stepTestStart,proto3" json:"stepTestStart,omitempty"`
	StepTestEnd          int32  `protobuf:"varint,7,opt,name=stepTestEnd,proto3" json:"stepTestEnd,omitempty"`
	StepTestStep         int32  `protobuf:"varint,8,opt,name=stepTestStep,proto3" json:"stepTestStep,omitempty"`
	StepTestSleepSeconds int32  `protobuf:"varint,9,opt,name=stepTestSleepSeconds,proto3" json:"stepTestSleepSeconds,omitempty"`
}

func (x *TestRequest) Reset() {
	*x = TestRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_stresstestloader_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRequest) ProtoMessage() {}

func (x *TestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_stresstestloader_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRequest.ProtoReflect.Descriptor instead.
func (*TestRequest) Descriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{1}
}

func (x *TestRequest) GetLoadtestExec() string {
	if x != nil {
		return x.LoadtestExec
	}
	return ""
}

func (x *TestRequest) GetS3() string {
	if x != nil {
		return x.S3
	}
	return ""
}

func (x *TestRequest) GetEnvVariableList() []*EnvVariable {
	if x != nil {
		return x.EnvVariableList
	}
	return nil
}

func (x *TestRequest) GetNumberOfGames() int32 {
	if x != nil {
		return x.NumberOfGames
	}
	return 0
}

func (x *TestRequest) GetS3Key() string {
	if x != nil {
		return x.S3Key
	}
	return ""
}

func (x *TestRequest) GetStepTestStart() int32 {
	if x != nil {
		return x.StepTestStart
	}
	return 0
}

func (x *TestRequest) GetStepTestEnd() int32 {
	if x != nil {
		return x.StepTestEnd
	}
	return 0
}

func (x *TestRequest) GetStepTestStep() int32 {
	if x != nil {
		return x.StepTestStep
	}
	return 0
}

func (x *TestRequest) GetStepTestSleepSeconds() int32 {
	if x != nil {
		return x.StepTestSleepSeconds
	}
	return 0
}

type StressTestConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkingFolder string                           `protobuf:"bytes,1,opt,name=workingFolder,proto3" json:"workingFolder,omitempty"`
	DebugL        *StressTestConfig_DebugLevelEnum `protobuf:"varint,4,opt,name=debugL,proto3,enum=stresstestloader.StressTestConfig_DebugLevelEnum,oneof" json:"debugL,omitempty"`
	ListenPort    int32                            `protobuf:"varint,5,opt,name=listenPort,proto3" json:"listenPort,omitempty"`
}

func (x *StressTestConfig) Reset() {
	*x = StressTestConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_stresstestloader_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StressTestConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StressTestConfig) ProtoMessage() {}

func (x *StressTestConfig) ProtoReflect() protoreflect.Message {
	mi := &file_proto_stresstestloader_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StressTestConfig.ProtoReflect.Descriptor instead.
func (*StressTestConfig) Descriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{2}
}

func (x *StressTestConfig) GetWorkingFolder() string {
	if x != nil {
		return x.WorkingFolder
	}
	return ""
}

func (x *StressTestConfig) GetDebugL() StressTestConfig_DebugLevelEnum {
	if x != nil && x.DebugL != nil {
		return *x.DebugL
	}
	return StressTestConfig_PanicLevel
}

func (x *StressTestConfig) GetListenPort() int32 {
	if x != nil {
		return x.ListenPort
	}
	return 0
}

// The response for the StartStressTest
type TestReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *TestReply) Reset() {
	*x = TestReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_stresstestloader_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestReply) ProtoMessage() {}

func (x *TestReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_stresstestloader_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestReply.ProtoReflect.Descriptor instead.
func (*TestReply) Descriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{3}
}

func (x *TestReply) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_proto_stresstestloader_proto protoreflect.FileDescriptor

var file_proto_stresstestloader_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65,
	0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10,
	0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72,
	0x22, 0x43, 0x0a, 0x0b, 0x45, 0x6e, 0x76, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x45, 0x6e, 0x76, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x45, 0x6e, 0x76, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x6e, 0x76,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x6e, 0x76,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xe6, 0x02, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x6f, 0x61, 0x64, 0x74, 0x65, 0x73,
	0x74, 0x45, 0x78, 0x65, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6c, 0x6f, 0x61,
	0x64, 0x74, 0x65, 0x73, 0x74, 0x45, 0x78, 0x65, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x33, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x73, 0x33, 0x12, 0x47, 0x0a, 0x0f, 0x65, 0x6e, 0x76,
	0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x6c,
	0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x45, 0x6e, 0x76, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c,
	0x65, 0x52, 0x0f, 0x65, 0x6e, 0x76, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x4f, 0x66, 0x47, 0x61,
	0x6d, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x4f, 0x66, 0x47, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x33, 0x6b, 0x65,
	0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x33, 0x6b, 0x65, 0x79, 0x12, 0x24,
	0x0a, 0x0d, 0x73, 0x74, 0x65, 0x70, 0x54, 0x65, 0x73, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x73, 0x74, 0x65, 0x70, 0x54, 0x65, 0x73, 0x74, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x74, 0x65, 0x70, 0x54, 0x65, 0x73, 0x74,
	0x45, 0x6e, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x73, 0x74, 0x65, 0x70, 0x54,
	0x65, 0x73, 0x74, 0x45, 0x6e, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x74, 0x65, 0x70, 0x54, 0x65,
	0x73, 0x74, 0x53, 0x74, 0x65, 0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x73, 0x74,
	0x65, 0x70, 0x54, 0x65, 0x73, 0x74, 0x53, 0x74, 0x65, 0x70, 0x12, 0x32, 0x0a, 0x14, 0x73, 0x74,
	0x65, 0x70, 0x54, 0x65, 0x73, 0x74, 0x53, 0x6c, 0x65, 0x65, 0x70, 0x53, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x14, 0x73, 0x74, 0x65, 0x70, 0x54, 0x65,
	0x73, 0x74, 0x53, 0x6c, 0x65, 0x65, 0x70, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0xb3,
	0x02, 0x0a, 0x10, 0x53, 0x74, 0x72, 0x65, 0x73, 0x73, 0x54, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x24, 0x0a, 0x0d, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x46, 0x6f,
	0x6c, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x77, 0x6f, 0x72, 0x6b,
	0x69, 0x6e, 0x67, 0x46, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x4e, 0x0a, 0x06, 0x64, 0x65, 0x62,
	0x75, 0x67, 0x4c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x31, 0x2e, 0x73, 0x74, 0x72, 0x65,
	0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x72,
	0x65, 0x73, 0x73, 0x54, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x44, 0x65,
	0x62, 0x75, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x45, 0x6e, 0x75, 0x6d, 0x48, 0x00, 0x52, 0x06,
	0x64, 0x65, 0x62, 0x75, 0x67, 0x4c, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x69, 0x73,
	0x74, 0x65, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x6c,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x7e, 0x0a, 0x0e, 0x44, 0x65, 0x62,
	0x75, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x0e, 0x0a, 0x0a, 0x50,
	0x61, 0x6e, 0x69, 0x63, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x46,
	0x61, 0x74, 0x61, 0x6c, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x57,
	0x61, 0x72, 0x6e, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x49, 0x6e,
	0x66, 0x6f, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x65, 0x62,
	0x75, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x72, 0x61,
	0x63, 0x65, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x06, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x64, 0x65,
	0x62, 0x75, 0x67, 0x4c, 0x22, 0x23, 0x0a, 0x09, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x63, 0x0a, 0x10, 0x53, 0x74, 0x72,
	0x65, 0x73, 0x73, 0x54, 0x65, 0x73, 0x74, 0x4c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x12, 0x4f, 0x0a,
	0x0f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x74, 0x72, 0x65, 0x73, 0x73, 0x54, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x6c, 0x6f, 0x61,
	0x64, 0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64,
	0x65, 0x72, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x49,
	0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x64, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73,
	0x2d, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x2d, 0x74, 0x65, 0x73,
	0x74, 0x2d, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x74,
	0x65, 0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_stresstestloader_proto_rawDescOnce sync.Once
	file_proto_stresstestloader_proto_rawDescData = file_proto_stresstestloader_proto_rawDesc
)

func file_proto_stresstestloader_proto_rawDescGZIP() []byte {
	file_proto_stresstestloader_proto_rawDescOnce.Do(func() {
		file_proto_stresstestloader_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_stresstestloader_proto_rawDescData)
	})
	return file_proto_stresstestloader_proto_rawDescData
}

var file_proto_stresstestloader_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_stresstestloader_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_stresstestloader_proto_goTypes = []interface{}{
	(StressTestConfig_DebugLevelEnum)(0), // 0: stresstestloader.StressTestConfig.DebugLevelEnum
	(*EnvVariable)(nil),                  // 1: stresstestloader.EnvVariable
	(*TestRequest)(nil),                  // 2: stresstestloader.TestRequest
	(*StressTestConfig)(nil),             // 3: stresstestloader.StressTestConfig
	(*TestReply)(nil),                    // 4: stresstestloader.TestReply
}
var file_proto_stresstestloader_proto_depIdxs = []int32{
	1, // 0: stresstestloader.TestRequest.envVariableList:type_name -> stresstestloader.EnvVariable
	0, // 1: stresstestloader.StressTestConfig.debugL:type_name -> stresstestloader.StressTestConfig.DebugLevelEnum
	2, // 2: stresstestloader.StressTestLoader.StartStressTest:input_type -> stresstestloader.TestRequest
	4, // 3: stresstestloader.StressTestLoader.StartStressTest:output_type -> stresstestloader.TestReply
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_stresstestloader_proto_init() }
func file_proto_stresstestloader_proto_init() {
	if File_proto_stresstestloader_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_stresstestloader_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvVariable); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_stresstestloader_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_stresstestloader_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StressTestConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_stresstestloader_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_stresstestloader_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_stresstestloader_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_stresstestloader_proto_goTypes,
		DependencyIndexes: file_proto_stresstestloader_proto_depIdxs,
		EnumInfos:         file_proto_stresstestloader_proto_enumTypes,
		MessageInfos:      file_proto_stresstestloader_proto_msgTypes,
	}.Build()
	File_proto_stresstestloader_proto = out.File
	file_proto_stresstestloader_proto_rawDesc = nil
	file_proto_stresstestloader_proto_goTypes = nil
	file_proto_stresstestloader_proto_depIdxs = nil
}
