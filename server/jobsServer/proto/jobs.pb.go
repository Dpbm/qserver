// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        v5.29.3
// source: jobs.proto

package jobs_server

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type JobProperties struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	ResultTypeCounts    bool                   `protobuf:"varint,1,opt,name=resultTypeCounts,proto3" json:"resultTypeCounts,omitempty"`
	ResultTypeQuasiDist bool                   `protobuf:"varint,2,opt,name=resultTypeQuasiDist,proto3" json:"resultTypeQuasiDist,omitempty"`
	ResultTypeExpVal    bool                   `protobuf:"varint,3,opt,name=resultTypeExpVal,proto3" json:"resultTypeExpVal,omitempty"`
	TargetSimulator     string                 `protobuf:"bytes,4,opt,name=targetSimulator,proto3" json:"targetSimulator,omitempty"`
	Metadata            string                 `protobuf:"bytes,5,opt,name=metadata,proto3" json:"metadata,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *JobProperties) Reset() {
	*x = JobProperties{}
	mi := &file_jobs_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobProperties) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobProperties) ProtoMessage() {}

func (x *JobProperties) ProtoReflect() protoreflect.Message {
	mi := &file_jobs_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobProperties.ProtoReflect.Descriptor instead.
func (*JobProperties) Descriptor() ([]byte, []int) {
	return file_jobs_proto_rawDescGZIP(), []int{0}
}

func (x *JobProperties) GetResultTypeCounts() bool {
	if x != nil {
		return x.ResultTypeCounts
	}
	return false
}

func (x *JobProperties) GetResultTypeQuasiDist() bool {
	if x != nil {
		return x.ResultTypeQuasiDist
	}
	return false
}

func (x *JobProperties) GetResultTypeExpVal() bool {
	if x != nil {
		return x.ResultTypeExpVal
	}
	return false
}

func (x *JobProperties) GetTargetSimulator() string {
	if x != nil {
		return x.TargetSimulator
	}
	return ""
}

func (x *JobProperties) GetMetadata() string {
	if x != nil {
		return x.Metadata
	}
	return ""
}

type JobData struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to DataStream:
	//
	//	*JobData_Properties
	//	*JobData_QasmChunk
	DataStream    isJobData_DataStream `protobuf_oneof:"DataStream"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobData) Reset() {
	*x = JobData{}
	mi := &file_jobs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobData) ProtoMessage() {}

func (x *JobData) ProtoReflect() protoreflect.Message {
	mi := &file_jobs_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobData.ProtoReflect.Descriptor instead.
func (*JobData) Descriptor() ([]byte, []int) {
	return file_jobs_proto_rawDescGZIP(), []int{1}
}

func (x *JobData) GetDataStream() isJobData_DataStream {
	if x != nil {
		return x.DataStream
	}
	return nil
}

func (x *JobData) GetProperties() *JobProperties {
	if x != nil {
		if x, ok := x.DataStream.(*JobData_Properties); ok {
			return x.Properties
		}
	}
	return nil
}

func (x *JobData) GetQasmChunk() string {
	if x != nil {
		if x, ok := x.DataStream.(*JobData_QasmChunk); ok {
			return x.QasmChunk
		}
	}
	return ""
}

type isJobData_DataStream interface {
	isJobData_DataStream()
}

type JobData_Properties struct {
	Properties *JobProperties `protobuf:"bytes,1,opt,name=properties,proto3,oneof"`
}

type JobData_QasmChunk struct {
	QasmChunk string `protobuf:"bytes,2,opt,name=qasmChunk,proto3,oneof"`
}

func (*JobData_Properties) isJobData_DataStream() {}

func (*JobData_QasmChunk) isJobData_DataStream() {}

type PendingJob struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PendingJob) Reset() {
	*x = PendingJob{}
	mi := &file_jobs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PendingJob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PendingJob) ProtoMessage() {}

func (x *PendingJob) ProtoReflect() protoreflect.Message {
	mi := &file_jobs_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PendingJob.ProtoReflect.Descriptor instead.
func (*PendingJob) Descriptor() ([]byte, []int) {
	return file_jobs_proto_rawDescGZIP(), []int{2}
}

func (x *PendingJob) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_jobs_proto protoreflect.FileDescriptor

var file_jobs_proto_rawDesc = string([]byte{
	0x0a, 0x0a, 0x6a, 0x6f, 0x62, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdf, 0x01, 0x0a,
	0x0d, 0x4a, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x12, 0x2a,
	0x0a, 0x10, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x79, 0x70, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12, 0x30, 0x0a, 0x13, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x54, 0x79, 0x70, 0x65, 0x51, 0x75, 0x61, 0x73, 0x69, 0x44, 0x69, 0x73,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x51, 0x75, 0x61, 0x73, 0x69, 0x44, 0x69, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x10,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x79, 0x70, 0x65, 0x45, 0x78, 0x70, 0x56, 0x61, 0x6c,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x45, 0x78, 0x70, 0x56, 0x61, 0x6c, 0x12, 0x28, 0x0a, 0x0f, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74,
	0x6f, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x69,
	0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x44, 0x61, 0x74, 0x61, 0x12, 0x30, 0x0a, 0x0a, 0x70, 0x72, 0x6f,
	0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x4a, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x48, 0x00, 0x52,
	0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x09, 0x71,
	0x61, 0x73, 0x6d, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x09, 0x71, 0x61, 0x73, 0x6d, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x42, 0x0c, 0x0a, 0x0a, 0x44,
	0x61, 0x74, 0x61, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x1c, 0x0a, 0x0a, 0x50, 0x65, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x4a, 0x6f, 0x62, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x32, 0x2b, 0x0a, 0x04, 0x4a, 0x6f, 0x62, 0x73, 0x12,
	0x23, 0x0a, 0x06, 0x41, 0x64, 0x64, 0x4a, 0x6f, 0x62, 0x12, 0x08, 0x2e, 0x4a, 0x6f, 0x62, 0x44,
	0x61, 0x74, 0x61, 0x1a, 0x0b, 0x2e, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x4a, 0x6f, 0x62,
	0x22, 0x00, 0x28, 0x01, 0x42, 0x1d, 0x5a, 0x1b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x44, 0x70, 0x62, 0x6d, 0x2f, 0x6a, 0x6f, 0x62, 0x73, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_jobs_proto_rawDescOnce sync.Once
	file_jobs_proto_rawDescData []byte
)

func file_jobs_proto_rawDescGZIP() []byte {
	file_jobs_proto_rawDescOnce.Do(func() {
		file_jobs_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_jobs_proto_rawDesc), len(file_jobs_proto_rawDesc)))
	})
	return file_jobs_proto_rawDescData
}

var file_jobs_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_jobs_proto_goTypes = []any{
	(*JobProperties)(nil), // 0: JobProperties
	(*JobData)(nil),       // 1: JobData
	(*PendingJob)(nil),    // 2: PendingJob
}
var file_jobs_proto_depIdxs = []int32{
	0, // 0: JobData.properties:type_name -> JobProperties
	1, // 1: Jobs.AddJob:input_type -> JobData
	2, // 2: Jobs.AddJob:output_type -> PendingJob
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_jobs_proto_init() }
func file_jobs_proto_init() {
	if File_jobs_proto != nil {
		return
	}
	file_jobs_proto_msgTypes[1].OneofWrappers = []any{
		(*JobData_Properties)(nil),
		(*JobData_QasmChunk)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_jobs_proto_rawDesc), len(file_jobs_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_jobs_proto_goTypes,
		DependencyIndexes: file_jobs_proto_depIdxs,
		MessageInfos:      file_jobs_proto_msgTypes,
	}.Build()
	File_jobs_proto = out.File
	file_jobs_proto_goTypes = nil
	file_jobs_proto_depIdxs = nil
}
