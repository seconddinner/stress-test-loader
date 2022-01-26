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

type MonitorServerConfig_DebugLevelEnum int32

const (
	MonitorServerConfig_PanicLevel MonitorServerConfig_DebugLevelEnum = 0
	MonitorServerConfig_FatalLevel MonitorServerConfig_DebugLevelEnum = 1
	MonitorServerConfig_ErrorLevel MonitorServerConfig_DebugLevelEnum = 2
	MonitorServerConfig_WarnLevel  MonitorServerConfig_DebugLevelEnum = 3
	MonitorServerConfig_InfoLevel  MonitorServerConfig_DebugLevelEnum = 4
	MonitorServerConfig_DebugLevel MonitorServerConfig_DebugLevelEnum = 5
	MonitorServerConfig_TraceLevel MonitorServerConfig_DebugLevelEnum = 6
)

// Enum value maps for MonitorServerConfig_DebugLevelEnum.
var (
	MonitorServerConfig_DebugLevelEnum_name = map[int32]string{
		0: "PanicLevel",
		1: "FatalLevel",
		2: "ErrorLevel",
		3: "WarnLevel",
		4: "InfoLevel",
		5: "DebugLevel",
		6: "TraceLevel",
	}
	MonitorServerConfig_DebugLevelEnum_value = map[string]int32{
		"PanicLevel": 0,
		"FatalLevel": 1,
		"ErrorLevel": 2,
		"WarnLevel":  3,
		"InfoLevel":  4,
		"DebugLevel": 5,
		"TraceLevel": 6,
	}
)

func (x MonitorServerConfig_DebugLevelEnum) Enum() *MonitorServerConfig_DebugLevelEnum {
	p := new(MonitorServerConfig_DebugLevelEnum)
	*p = x
	return p
}

func (x MonitorServerConfig_DebugLevelEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MonitorServerConfig_DebugLevelEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_stresstestloader_proto_enumTypes[0].Descriptor()
}

func (MonitorServerConfig_DebugLevelEnum) Type() protoreflect.EnumType {
	return &file_proto_stresstestloader_proto_enumTypes[0]
}

func (x MonitorServerConfig_DebugLevelEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MonitorServerConfig_DebugLevelEnum.Descriptor instead.
func (MonitorServerConfig_DebugLevelEnum) EnumDescriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{1, 0}
}

// The request to verify
type VerifyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proxyurl      *string `protobuf:"bytes,1,opt,name=proxyurl,proto3,oneof" json:"proxyurl,omitempty"`
	Url           string  `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Timeout       int32   `protobuf:"varint,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	PostUser      *string `protobuf:"bytes,4,opt,name=postUser,proto3,oneof" json:"postUser,omitempty"`
	PostPassword  *string `protobuf:"bytes,5,opt,name=postPassword,proto3,oneof" json:"postPassword,omitempty"`
	UrlParameters *string `protobuf:"bytes,6,opt,name=urlParameters,proto3,oneof" json:"urlParameters,omitempty"` // refer to https://docs.google.com/document/d/1sKtA2ExVS5gjjqGkGOldib8UO_qvMzWH8qNloCmuJI4/edit#
}

func (x *VerifyRequest) Reset() {
	*x = VerifyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_stresstestloader_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyRequest) ProtoMessage() {}

func (x *VerifyRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use VerifyRequest.ProtoReflect.Descriptor instead.
func (*VerifyRequest) Descriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{0}
}

func (x *VerifyRequest) GetProxyurl() string {
	if x != nil && x.Proxyurl != nil {
		return *x.Proxyurl
	}
	return ""
}

func (x *VerifyRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *VerifyRequest) GetTimeout() int32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *VerifyRequest) GetPostUser() string {
	if x != nil && x.PostUser != nil {
		return *x.PostUser
	}
	return ""
}

func (x *VerifyRequest) GetPostPassword() string {
	if x != nil && x.PostPassword != nil {
		return *x.PostPassword
	}
	return ""
}

func (x *VerifyRequest) GetUrlParameters() string {
	if x != nil && x.UrlParameters != nil {
		return *x.UrlParameters
	}
	return ""
}

type MonitorServerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VRList                   []*VerifyRequest                    `protobuf:"bytes,1,rep,name=VRList,proto3" json:"VRList,omitempty"`
	WaitSeconds              int32                               `protobuf:"varint,2,opt,name=waitSeconds,proto3" json:"waitSeconds,omitempty"`
	InitialWaitForProxyStart *int32                              `protobuf:"varint,3,opt,name=initialWaitForProxyStart,proto3,oneof" json:"initialWaitForProxyStart,omitempty"`
	DebugL                   *MonitorServerConfig_DebugLevelEnum `protobuf:"varint,4,opt,name=debugL,proto3,enum=stresstestloader.MonitorServerConfig_DebugLevelEnum,oneof" json:"debugL,omitempty"`
	ListenPort               int32                               `protobuf:"varint,5,opt,name=listenPort,proto3" json:"listenPort,omitempty"`
	ShutdownThreshHold       *int32                              `protobuf:"varint,6,opt,name=shutdownThreshHold,proto3,oneof" json:"shutdownThreshHold,omitempty"` // number of times we retry http before we shutdown squid
	ShutdownFlag             *bool                               `protobuf:"varint,7,opt,name=shutdownFlag,proto3,oneof" json:"shutdownFlag,omitempty"`
}

func (x *MonitorServerConfig) Reset() {
	*x = MonitorServerConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_stresstestloader_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MonitorServerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MonitorServerConfig) ProtoMessage() {}

func (x *MonitorServerConfig) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use MonitorServerConfig.ProtoReflect.Descriptor instead.
func (*MonitorServerConfig) Descriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{1}
}

func (x *MonitorServerConfig) GetVRList() []*VerifyRequest {
	if x != nil {
		return x.VRList
	}
	return nil
}

func (x *MonitorServerConfig) GetWaitSeconds() int32 {
	if x != nil {
		return x.WaitSeconds
	}
	return 0
}

func (x *MonitorServerConfig) GetInitialWaitForProxyStart() int32 {
	if x != nil && x.InitialWaitForProxyStart != nil {
		return *x.InitialWaitForProxyStart
	}
	return 0
}

func (x *MonitorServerConfig) GetDebugL() MonitorServerConfig_DebugLevelEnum {
	if x != nil && x.DebugL != nil {
		return *x.DebugL
	}
	return MonitorServerConfig_PanicLevel
}

func (x *MonitorServerConfig) GetListenPort() int32 {
	if x != nil {
		return x.ListenPort
	}
	return 0
}

func (x *MonitorServerConfig) GetShutdownThreshHold() int32 {
	if x != nil && x.ShutdownThreshHold != nil {
		return *x.ShutdownThreshHold
	}
	return 0
}

func (x *MonitorServerConfig) GetShutdownFlag() bool {
	if x != nil && x.ShutdownFlag != nil {
		return *x.ShutdownFlag
	}
	return false
}

// The response for the verify
type VerifyReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *VerifyReply) Reset() {
	*x = VerifyReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_stresstestloader_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyReply) ProtoMessage() {}

func (x *VerifyReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use VerifyReply.ProtoReflect.Descriptor instead.
func (*VerifyReply) Descriptor() ([]byte, []int) {
	return file_proto_stresstestloader_proto_rawDescGZIP(), []int{2}
}

func (x *VerifyReply) GetStatus() string {
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
	0x22, 0x8e, 0x02, 0x0a, 0x0d, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x75, 0x72, 0x6c,
	0x88, 0x01, 0x01, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12,
	0x1f, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x01, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x88, 0x01, 0x01,
	0x12, 0x27, 0x0a, 0x0c, 0x70, 0x6f, 0x73, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x0c, 0x70, 0x6f, 0x73, 0x74, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x88, 0x01, 0x01, 0x12, 0x29, 0x0a, 0x0d, 0x75, 0x72, 0x6c,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x03, 0x52, 0x0d, 0x75, 0x72, 0x6c, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x73, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x75, 0x72,
	0x6c, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x42, 0x0f,
	0x0a, 0x0d, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x42,
	0x10, 0x0a, 0x0e, 0x5f, 0x75, 0x72, 0x6c, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x73, 0x22, 0xd2, 0x04, 0x0a, 0x13, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x37, 0x0a, 0x06, 0x56, 0x52, 0x4c,
	0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x73, 0x74, 0x72, 0x65,
	0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x06, 0x56, 0x52, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x77, 0x61, 0x69, 0x74, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x77, 0x61, 0x69, 0x74, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x73, 0x12, 0x3f, 0x0a, 0x18, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x57,
	0x61, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x18, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61,
	0x6c, 0x57, 0x61, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x88, 0x01, 0x01, 0x12, 0x51, 0x0a, 0x06, 0x64, 0x65, 0x62, 0x75, 0x67, 0x4c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x34, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65,
	0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x44, 0x65, 0x62,
	0x75, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x45, 0x6e, 0x75, 0x6d, 0x48, 0x01, 0x52, 0x06, 0x64,
	0x65, 0x62, 0x75, 0x67, 0x4c, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x69, 0x73, 0x74,
	0x65, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x6c, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x33, 0x0a, 0x12, 0x73, 0x68, 0x75, 0x74,
	0x64, 0x6f, 0x77, 0x6e, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x48, 0x6f, 0x6c, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x02, 0x52, 0x12, 0x73, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e,
	0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x48, 0x6f, 0x6c, 0x64, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a,
	0x0c, 0x73, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x46, 0x6c, 0x61, 0x67, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x08, 0x48, 0x03, 0x52, 0x0c, 0x73, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x46,
	0x6c, 0x61, 0x67, 0x88, 0x01, 0x01, 0x22, 0x7e, 0x0a, 0x0e, 0x44, 0x65, 0x62, 0x75, 0x67, 0x4c,
	0x65, 0x76, 0x65, 0x6c, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x61, 0x6e, 0x69,
	0x63, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x46, 0x61, 0x74, 0x61,
	0x6c, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x57, 0x61, 0x72, 0x6e,
	0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x49, 0x6e, 0x66, 0x6f, 0x4c,
	0x65, 0x76, 0x65, 0x6c, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x65, 0x62, 0x75, 0x67, 0x4c,
	0x65, 0x76, 0x65, 0x6c, 0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x72, 0x61, 0x63, 0x65, 0x4c,
	0x65, 0x76, 0x65, 0x6c, 0x10, 0x06, 0x42, 0x1b, 0x0a, 0x19, 0x5f, 0x69, 0x6e, 0x69, 0x74, 0x69,
	0x61, 0x6c, 0x57, 0x61, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x4c, 0x42, 0x15,
	0x0a, 0x13, 0x5f, 0x73, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x54, 0x68, 0x72, 0x65, 0x73,
	0x68, 0x48, 0x6f, 0x6c, 0x64, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x73, 0x68, 0x75, 0x74, 0x64, 0x6f,
	0x77, 0x6e, 0x46, 0x6c, 0x61, 0x67, 0x22, 0x25, 0x0a, 0x0b, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x5c, 0x0a,
	0x08, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x50, 0x0a, 0x06, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x79, 0x12, 0x25, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74,
	0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x1d, 0x2e, 0x73, 0x74, 0x72,
	0x65, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x51, 0x5a, 0x4f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x64, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x2f, 0x6a, 0x78, 0x69, 0x65, 0x2d, 0x73, 0x64, 0x2f, 0x73,
	0x74, 0x72, 0x65, 0x73, 0x73, 0x2d, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x73,
	0x73, 0x2d, 0x74, 0x65, 0x73, 0x74, 0x2d, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2f, 0x73, 0x74,
	0x72, 0x65, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
var file_proto_stresstestloader_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_stresstestloader_proto_goTypes = []interface{}{
	(MonitorServerConfig_DebugLevelEnum)(0), // 0: stresstestloader.MonitorServerConfig.DebugLevelEnum
	(*VerifyRequest)(nil),                   // 1: stresstestloader.VerifyRequest
	(*MonitorServerConfig)(nil),             // 2: stresstestloader.MonitorServerConfig
	(*VerifyReply)(nil),                     // 3: stresstestloader.VerifyReply
}
var file_proto_stresstestloader_proto_depIdxs = []int32{
	1, // 0: stresstestloader.MonitorServerConfig.VRList:type_name -> stresstestloader.VerifyRequest
	0, // 1: stresstestloader.MonitorServerConfig.debugL:type_name -> stresstestloader.MonitorServerConfig.DebugLevelEnum
	2, // 2: stresstestloader.Verifier.Verify:input_type -> stresstestloader.MonitorServerConfig
	3, // 3: stresstestloader.Verifier.Verify:output_type -> stresstestloader.VerifyReply
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
			switch v := v.(*VerifyRequest); i {
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
			switch v := v.(*MonitorServerConfig); i {
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
			switch v := v.(*VerifyReply); i {
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
	file_proto_stresstestloader_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_proto_stresstestloader_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_stresstestloader_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
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
