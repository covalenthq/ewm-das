// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v5.28.0--dev
// source: proto/workload.proto

package workloadpb

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

type Workload struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	ChainId             uint64                 `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	BlockHeight         uint64                 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	BlobIndex           uint64                 `protobuf:"varint,3,opt,name=blob_index,json=blobIndex,proto3" json:"blob_index,omitempty"`
	ExpirationTimestamp uint64                 `protobuf:"varint,4,opt,name=expiration_timestamp,json=expirationTimestamp,proto3" json:"expiration_timestamp,omitempty"`
	Hash                []byte                 `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
	BlockHash           []byte                 `protobuf:"bytes,6,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	SpecimenHash        []byte                 `protobuf:"bytes,7,opt,name=specimen_hash,json=specimenHash,proto3" json:"specimen_hash,omitempty"`
	Commitment          []byte                 `protobuf:"bytes,8,opt,name=commitment,proto3" json:"commitment,omitempty"`
	IpfsCid             []byte                 `protobuf:"bytes,9,opt,name=ipfs_cid,json=ipfsCid,proto3" json:"ipfs_cid,omitempty"`
	Challenge           []byte                 `protobuf:"bytes,10,opt,name=challenge,proto3" json:"challenge,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *Workload) Reset() {
	*x = Workload{}
	mi := &file_proto_workload_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Workload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workload) ProtoMessage() {}

func (x *Workload) ProtoReflect() protoreflect.Message {
	mi := &file_proto_workload_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Workload.ProtoReflect.Descriptor instead.
func (*Workload) Descriptor() ([]byte, []int) {
	return file_proto_workload_proto_rawDescGZIP(), []int{0}
}

func (x *Workload) GetChainId() uint64 {
	if x != nil {
		return x.ChainId
	}
	return 0
}

func (x *Workload) GetBlockHeight() uint64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *Workload) GetBlobIndex() uint64 {
	if x != nil {
		return x.BlobIndex
	}
	return 0
}

func (x *Workload) GetExpirationTimestamp() uint64 {
	if x != nil {
		return x.ExpirationTimestamp
	}
	return 0
}

func (x *Workload) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Workload) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *Workload) GetSpecimenHash() []byte {
	if x != nil {
		return x.SpecimenHash
	}
	return nil
}

func (x *Workload) GetCommitment() []byte {
	if x != nil {
		return x.Commitment
	}
	return nil
}

func (x *Workload) GetIpfsCid() []byte {
	if x != nil {
		return x.IpfsCid
	}
	return nil
}

func (x *Workload) GetChallenge() []byte {
	if x != nil {
		return x.Challenge
	}
	return nil
}

type SignedWorkload struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Workload      *Workload              `protobuf:"bytes,1,opt,name=workload,proto3" json:"workload,omitempty"`
	Signature     []byte                 `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SignedWorkload) Reset() {
	*x = SignedWorkload{}
	mi := &file_proto_workload_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignedWorkload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedWorkload) ProtoMessage() {}

func (x *SignedWorkload) ProtoReflect() protoreflect.Message {
	mi := &file_proto_workload_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedWorkload.ProtoReflect.Descriptor instead.
func (*SignedWorkload) Descriptor() ([]byte, []int) {
	return file_proto_workload_proto_rawDescGZIP(), []int{1}
}

func (x *SignedWorkload) GetWorkload() *Workload {
	if x != nil {
		return x.Workload
	}
	return nil
}

func (x *SignedWorkload) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type SignedWorkloadCollection struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Workloads           []*SignedWorkload      `protobuf:"bytes,1,rep,name=workloads,proto3" json:"workloads,omitempty"`
	NextUpdateTimestamp uint64                 `protobuf:"varint,2,opt,name=next_update_timestamp,json=nextUpdateTimestamp,proto3" json:"next_update_timestamp,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *SignedWorkloadCollection) Reset() {
	*x = SignedWorkloadCollection{}
	mi := &file_proto_workload_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignedWorkloadCollection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignedWorkloadCollection) ProtoMessage() {}

func (x *SignedWorkloadCollection) ProtoReflect() protoreflect.Message {
	mi := &file_proto_workload_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignedWorkloadCollection.ProtoReflect.Descriptor instead.
func (*SignedWorkloadCollection) Descriptor() ([]byte, []int) {
	return file_proto_workload_proto_rawDescGZIP(), []int{2}
}

func (x *SignedWorkloadCollection) GetWorkloads() []*SignedWorkload {
	if x != nil {
		return x.Workloads
	}
	return nil
}

func (x *SignedWorkloadCollection) GetNextUpdateTimestamp() uint64 {
	if x != nil {
		return x.NextUpdateTimestamp
	}
	return 0
}

type StoreRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Workload      *SignedWorkload        `protobuf:"bytes,1,opt,name=workload,proto3" json:"workload,omitempty"`
	Timestamp     uint64                 `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	CellIndex     uint64                 `protobuf:"varint,3,opt,name=cell_index,json=cellIndex,proto3" json:"cell_index,omitempty"`
	Proof         []byte                 `protobuf:"bytes,4,opt,name=proof,proto3" json:"proof,omitempty"`
	Cell          []byte                 `protobuf:"bytes,5,opt,name=cell,proto3" json:"cell,omitempty"`
	Version       string                 `protobuf:"bytes,6,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StoreRequest) Reset() {
	*x = StoreRequest{}
	mi := &file_proto_workload_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreRequest) ProtoMessage() {}

func (x *StoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_workload_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreRequest.ProtoReflect.Descriptor instead.
func (*StoreRequest) Descriptor() ([]byte, []int) {
	return file_proto_workload_proto_rawDescGZIP(), []int{3}
}

func (x *StoreRequest) GetWorkload() *SignedWorkload {
	if x != nil {
		return x.Workload
	}
	return nil
}

func (x *StoreRequest) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *StoreRequest) GetCellIndex() uint64 {
	if x != nil {
		return x.CellIndex
	}
	return 0
}

func (x *StoreRequest) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *StoreRequest) GetCell() []byte {
	if x != nil {
		return x.Cell
	}
	return nil
}

func (x *StoreRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_proto_workload_proto protoreflect.FileDescriptor

var file_proto_workload_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x22, 0xcb, 0x02, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x19, 0x0a,
	0x08, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x62,
	0x6c, 0x6f, 0x62, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x09, 0x62, 0x6c, 0x6f, 0x62, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x31, 0x0a, 0x14, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x13, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68,
	0x12, 0x23, 0x0a, 0x0d, 0x73, 0x70, 0x65, 0x63, 0x69, 0x6d, 0x65, 0x6e, 0x5f, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x73, 0x70, 0x65, 0x63, 0x69, 0x6d, 0x65,
	0x6e, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x6d,
	0x65, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x70, 0x66, 0x73, 0x5f, 0x63, 0x69,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x69, 0x70, 0x66, 0x73, 0x43, 0x69, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x22, 0x5e,
	0x0a, 0x0e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x2e, 0x0a, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x57, 0x6f,
	0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x86,
	0x01, 0x0a, 0x18, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x09, 0x77,
	0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f,
	0x61, 0x64, 0x73, 0x12, 0x32, 0x0a, 0x15, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x13, 0x6e, 0x65, 0x78, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0xc5, 0x01, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x08, 0x77, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x77, 0x6f, 0x72,
	0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x57, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x65, 0x6c, 0x6c, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x63, 0x65, 0x6c, 0x6c, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x70,
	0x72, 0x6f, 0x6f, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f,
	0x66, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x6c, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x63, 0x65, 0x6c, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_workload_proto_rawDescOnce sync.Once
	file_proto_workload_proto_rawDescData = file_proto_workload_proto_rawDesc
)

func file_proto_workload_proto_rawDescGZIP() []byte {
	file_proto_workload_proto_rawDescOnce.Do(func() {
		file_proto_workload_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_workload_proto_rawDescData)
	})
	return file_proto_workload_proto_rawDescData
}

var file_proto_workload_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_workload_proto_goTypes = []any{
	(*Workload)(nil),                 // 0: workload.Workload
	(*SignedWorkload)(nil),           // 1: workload.SignedWorkload
	(*SignedWorkloadCollection)(nil), // 2: workload.SignedWorkloadCollection
	(*StoreRequest)(nil),             // 3: workload.StoreRequest
}
var file_proto_workload_proto_depIdxs = []int32{
	0, // 0: workload.SignedWorkload.workload:type_name -> workload.Workload
	1, // 1: workload.SignedWorkloadCollection.workloads:type_name -> workload.SignedWorkload
	1, // 2: workload.StoreRequest.workload:type_name -> workload.SignedWorkload
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_workload_proto_init() }
func file_proto_workload_proto_init() {
	if File_proto_workload_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_workload_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_workload_proto_goTypes,
		DependencyIndexes: file_proto_workload_proto_depIdxs,
		MessageInfos:      file_proto_workload_proto_msgTypes,
	}.Build()
	File_proto_workload_proto = out.File
	file_proto_workload_proto_rawDesc = nil
	file_proto_workload_proto_goTypes = nil
	file_proto_workload_proto_depIdxs = nil
}
