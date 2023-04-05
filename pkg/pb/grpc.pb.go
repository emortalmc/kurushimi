// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: grpc.proto

package pb

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

type QueueByPlayerErrorResponse_ErrorReason int32

const (
	QueueByPlayerErrorResponse_ALREADY_IN_QUEUE    QueueByPlayerErrorResponse_ErrorReason = 0
	QueueByPlayerErrorResponse_INVALID_GAME_MODE   QueueByPlayerErrorResponse_ErrorReason = 1
	QueueByPlayerErrorResponse_GAME_MODE_DISABLED  QueueByPlayerErrorResponse_ErrorReason = 2
	QueueByPlayerErrorResponse_INVALID_MAP         QueueByPlayerErrorResponse_ErrorReason = 3
	QueueByPlayerErrorResponse_PARTY_TOO_LARGE     QueueByPlayerErrorResponse_ErrorReason = 4
	QueueByPlayerErrorResponse_PARTIES_NOT_ALLOWED QueueByPlayerErrorResponse_ErrorReason = 5
)

// Enum value maps for QueueByPlayerErrorResponse_ErrorReason.
var (
	QueueByPlayerErrorResponse_ErrorReason_name = map[int32]string{
		0: "ALREADY_IN_QUEUE",
		1: "INVALID_GAME_MODE",
		2: "GAME_MODE_DISABLED",
		3: "INVALID_MAP",
		4: "PARTY_TOO_LARGE",
		5: "PARTIES_NOT_ALLOWED",
	}
	QueueByPlayerErrorResponse_ErrorReason_value = map[string]int32{
		"ALREADY_IN_QUEUE":    0,
		"INVALID_GAME_MODE":   1,
		"GAME_MODE_DISABLED":  2,
		"INVALID_MAP":         3,
		"PARTY_TOO_LARGE":     4,
		"PARTIES_NOT_ALLOWED": 5,
	}
)

func (x QueueByPlayerErrorResponse_ErrorReason) Enum() *QueueByPlayerErrorResponse_ErrorReason {
	p := new(QueueByPlayerErrorResponse_ErrorReason)
	*p = x
	return p
}

func (x QueueByPlayerErrorResponse_ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QueueByPlayerErrorResponse_ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_proto_enumTypes[0].Descriptor()
}

func (QueueByPlayerErrorResponse_ErrorReason) Type() protoreflect.EnumType {
	return &file_grpc_proto_enumTypes[0]
}

func (x QueueByPlayerErrorResponse_ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueueByPlayerErrorResponse_ErrorReason.Descriptor instead.
func (QueueByPlayerErrorResponse_ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{2, 0}
}

type DequeueByPlayerErrorResponse_ErrorReason int32

const (
	DequeueByPlayerErrorResponse_NOT_IN_QUEUE               DequeueByPlayerErrorResponse_ErrorReason = 0
	DequeueByPlayerErrorResponse_NO_PERMISSION              DequeueByPlayerErrorResponse_ErrorReason = 1
	DequeueByPlayerErrorResponse_ALREADY_MARKED_FOR_DEQUEUE DequeueByPlayerErrorResponse_ErrorReason = 2
)

// Enum value maps for DequeueByPlayerErrorResponse_ErrorReason.
var (
	DequeueByPlayerErrorResponse_ErrorReason_name = map[int32]string{
		0: "NOT_IN_QUEUE",
		1: "NO_PERMISSION",
		2: "ALREADY_MARKED_FOR_DEQUEUE",
	}
	DequeueByPlayerErrorResponse_ErrorReason_value = map[string]int32{
		"NOT_IN_QUEUE":               0,
		"NO_PERMISSION":              1,
		"ALREADY_MARKED_FOR_DEQUEUE": 2,
	}
)

func (x DequeueByPlayerErrorResponse_ErrorReason) Enum() *DequeueByPlayerErrorResponse_ErrorReason {
	p := new(DequeueByPlayerErrorResponse_ErrorReason)
	*p = x
	return p
}

func (x DequeueByPlayerErrorResponse_ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DequeueByPlayerErrorResponse_ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_proto_enumTypes[1].Descriptor()
}

func (DequeueByPlayerErrorResponse_ErrorReason) Type() protoreflect.EnumType {
	return &file_grpc_proto_enumTypes[1]
}

func (x DequeueByPlayerErrorResponse_ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DequeueByPlayerErrorResponse_ErrorReason.Descriptor instead.
func (DequeueByPlayerErrorResponse_ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{5, 0}
}

type ChangePlayerMapVoteErrorResponse_ErrorReason int32

const (
	ChangePlayerMapVoteErrorResponse_INVALID_MAP  ChangePlayerMapVoteErrorResponse_ErrorReason = 0
	ChangePlayerMapVoteErrorResponse_NOT_IN_QUEUE ChangePlayerMapVoteErrorResponse_ErrorReason = 1
)

// Enum value maps for ChangePlayerMapVoteErrorResponse_ErrorReason.
var (
	ChangePlayerMapVoteErrorResponse_ErrorReason_name = map[int32]string{
		0: "INVALID_MAP",
		1: "NOT_IN_QUEUE",
	}
	ChangePlayerMapVoteErrorResponse_ErrorReason_value = map[string]int32{
		"INVALID_MAP":  0,
		"NOT_IN_QUEUE": 1,
	}
)

func (x ChangePlayerMapVoteErrorResponse_ErrorReason) Enum() *ChangePlayerMapVoteErrorResponse_ErrorReason {
	p := new(ChangePlayerMapVoteErrorResponse_ErrorReason)
	*p = x
	return p
}

func (x ChangePlayerMapVoteErrorResponse_ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChangePlayerMapVoteErrorResponse_ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_proto_enumTypes[2].Descriptor()
}

func (ChangePlayerMapVoteErrorResponse_ErrorReason) Type() protoreflect.EnumType {
	return &file_grpc_proto_enumTypes[2]
}

func (x ChangePlayerMapVoteErrorResponse_ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChangePlayerMapVoteErrorResponse_ErrorReason.Descriptor instead.
func (ChangePlayerMapVoteErrorResponse_ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{8, 0}
}

type QueueByPlayerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameModeId string `protobuf:"bytes,1,opt,name=game_mode_id,json=gameModeId,proto3" json:"game_mode_id,omitempty"`
	// player_id of type UUID
	PlayerId string `protobuf:"bytes,2,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	// private_game defaults to false. If true, gamemode will be created as an instant private game.
	PrivateGame *bool `protobuf:"varint,3,opt,name=private_game,json=privateGame,proto3,oneof" json:"private_game,omitempty"`
	// voting
	MapId *string `protobuf:"bytes,4,opt,name=map_id,json=mapId,proto3,oneof" json:"map_id,omitempty"`
	// auto_teleport if true or not specified, the proxy will listen for match messages and teleport the player to the match.
	// if false, the proxy will not auto-teleport the player and you should listen for match messages yourself.
	//
	// e.g. this is used by the proxy for lobby matchmaking as it is handled differently
	// to when the player is already connected to the server.
	AutoTeleport *bool `protobuf:"varint,5,opt,name=auto_teleport,json=autoTeleport,proto3,oneof" json:"auto_teleport,omitempty"`
}

func (x *QueueByPlayerRequest) Reset() {
	*x = QueueByPlayerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueByPlayerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueByPlayerRequest) ProtoMessage() {}

func (x *QueueByPlayerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueByPlayerRequest.ProtoReflect.Descriptor instead.
func (*QueueByPlayerRequest) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{0}
}

func (x *QueueByPlayerRequest) GetGameModeId() string {
	if x != nil {
		return x.GameModeId
	}
	return ""
}

func (x *QueueByPlayerRequest) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *QueueByPlayerRequest) GetPrivateGame() bool {
	if x != nil && x.PrivateGame != nil {
		return *x.PrivateGame
	}
	return false
}

func (x *QueueByPlayerRequest) GetMapId() string {
	if x != nil && x.MapId != nil {
		return *x.MapId
	}
	return ""
}

func (x *QueueByPlayerRequest) GetAutoTeleport() bool {
	if x != nil && x.AutoTeleport != nil {
		return *x.AutoTeleport
	}
	return false
}

type QueueByPlayerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *QueueByPlayerResponse) Reset() {
	*x = QueueByPlayerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueByPlayerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueByPlayerResponse) ProtoMessage() {}

func (x *QueueByPlayerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueByPlayerResponse.ProtoReflect.Descriptor instead.
func (*QueueByPlayerResponse) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{1}
}

type QueueByPlayerErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reason QueueByPlayerErrorResponse_ErrorReason `protobuf:"varint,1,opt,name=reason,proto3,enum=dev.emortal.kurushimi.grpc.QueueByPlayerErrorResponse_ErrorReason" json:"reason,omitempty"`
}

func (x *QueueByPlayerErrorResponse) Reset() {
	*x = QueueByPlayerErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueByPlayerErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueByPlayerErrorResponse) ProtoMessage() {}

func (x *QueueByPlayerErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueByPlayerErrorResponse.ProtoReflect.Descriptor instead.
func (*QueueByPlayerErrorResponse) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{2}
}

func (x *QueueByPlayerErrorResponse) GetReason() QueueByPlayerErrorResponse_ErrorReason {
	if x != nil {
		return x.Reason
	}
	return QueueByPlayerErrorResponse_ALREADY_IN_QUEUE
}

type DequeueByPlayerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// player_id of type UUID
	PlayerId string `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
}

func (x *DequeueByPlayerRequest) Reset() {
	*x = DequeueByPlayerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DequeueByPlayerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DequeueByPlayerRequest) ProtoMessage() {}

func (x *DequeueByPlayerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DequeueByPlayerRequest.ProtoReflect.Descriptor instead.
func (*DequeueByPlayerRequest) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{3}
}

func (x *DequeueByPlayerRequest) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

type DequeueByPlayerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DequeueByPlayerResponse) Reset() {
	*x = DequeueByPlayerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DequeueByPlayerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DequeueByPlayerResponse) ProtoMessage() {}

func (x *DequeueByPlayerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DequeueByPlayerResponse.ProtoReflect.Descriptor instead.
func (*DequeueByPlayerResponse) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{4}
}

type DequeueByPlayerErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reason DequeueByPlayerErrorResponse_ErrorReason `protobuf:"varint,1,opt,name=reason,proto3,enum=dev.emortal.kurushimi.grpc.DequeueByPlayerErrorResponse_ErrorReason" json:"reason,omitempty"`
}

func (x *DequeueByPlayerErrorResponse) Reset() {
	*x = DequeueByPlayerErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DequeueByPlayerErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DequeueByPlayerErrorResponse) ProtoMessage() {}

func (x *DequeueByPlayerErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DequeueByPlayerErrorResponse.ProtoReflect.Descriptor instead.
func (*DequeueByPlayerErrorResponse) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{5}
}

func (x *DequeueByPlayerErrorResponse) GetReason() DequeueByPlayerErrorResponse_ErrorReason {
	if x != nil {
		return x.Reason
	}
	return DequeueByPlayerErrorResponse_NOT_IN_QUEUE
}

type ChangePlayerMapVoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId string `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	MapId    string `protobuf:"bytes,2,opt,name=map_id,json=mapId,proto3" json:"map_id,omitempty"`
}

func (x *ChangePlayerMapVoteRequest) Reset() {
	*x = ChangePlayerMapVoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangePlayerMapVoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangePlayerMapVoteRequest) ProtoMessage() {}

func (x *ChangePlayerMapVoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangePlayerMapVoteRequest.ProtoReflect.Descriptor instead.
func (*ChangePlayerMapVoteRequest) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{6}
}

func (x *ChangePlayerMapVoteRequest) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *ChangePlayerMapVoteRequest) GetMapId() string {
	if x != nil {
		return x.MapId
	}
	return ""
}

type ChangePlayerMapVoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ChangePlayerMapVoteResponse) Reset() {
	*x = ChangePlayerMapVoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangePlayerMapVoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangePlayerMapVoteResponse) ProtoMessage() {}

func (x *ChangePlayerMapVoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangePlayerMapVoteResponse.ProtoReflect.Descriptor instead.
func (*ChangePlayerMapVoteResponse) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{7}
}

type ChangePlayerMapVoteErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reason ChangePlayerMapVoteErrorResponse_ErrorReason `protobuf:"varint,1,opt,name=reason,proto3,enum=dev.emortal.kurushimi.grpc.ChangePlayerMapVoteErrorResponse_ErrorReason" json:"reason,omitempty"`
}

func (x *ChangePlayerMapVoteErrorResponse) Reset() {
	*x = ChangePlayerMapVoteErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangePlayerMapVoteErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangePlayerMapVoteErrorResponse) ProtoMessage() {}

func (x *ChangePlayerMapVoteErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangePlayerMapVoteErrorResponse.ProtoReflect.Descriptor instead.
func (*ChangePlayerMapVoteErrorResponse) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{8}
}

func (x *ChangePlayerMapVoteErrorResponse) GetReason() ChangePlayerMapVoteErrorResponse_ErrorReason {
	if x != nil {
		return x.Reason
	}
	return ChangePlayerMapVoteErrorResponse_INVALID_MAP
}

var File_grpc_proto protoreflect.FileDescriptor

var file_grpc_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x64, 0x65,
	0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x6b, 0x75, 0x72, 0x75, 0x73, 0x68,
	0x69, 0x6d, 0x69, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x22, 0xf1, 0x01, 0x0a, 0x14, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x20, 0x0a, 0x0c, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x67, 0x61, 0x6d, 0x65, 0x4d, 0x6f, 0x64,
	0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x26, 0x0a, 0x0c, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x67, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x0b, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x47, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1a, 0x0a, 0x06, 0x6d, 0x61, 0x70, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x05, 0x6d, 0x61, 0x70, 0x49,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x28, 0x0a, 0x0d, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x74, 0x65, 0x6c,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x48, 0x02, 0x52, 0x0c, 0x61,
	0x75, 0x74, 0x6f, 0x54, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x88, 0x01, 0x01, 0x42, 0x0f,
	0x0a, 0x0d, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x67, 0x61, 0x6d, 0x65, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x6d, 0x61, 0x70, 0x5f, 0x69, 0x64, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x61,
	0x75, 0x74, 0x6f, 0x5f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x17, 0x0a, 0x15,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x8c, 0x02, 0x0a, 0x1a, 0x51, 0x75, 0x65, 0x75, 0x65, 0x42,
	0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x42, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72, 0x74,
	0x61, 0x6c, 0x2e, 0x6b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x22, 0x91, 0x01, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x12, 0x14, 0x0a, 0x10, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x49, 0x4e, 0x5f, 0x51,
	0x55, 0x45, 0x55, 0x45, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49,
	0x44, 0x5f, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x10, 0x01, 0x12, 0x16, 0x0a,
	0x12, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x44, 0x49, 0x53, 0x41, 0x42,
	0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44,
	0x5f, 0x4d, 0x41, 0x50, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x41, 0x52, 0x54, 0x59, 0x5f,
	0x54, 0x4f, 0x4f, 0x5f, 0x4c, 0x41, 0x52, 0x47, 0x45, 0x10, 0x04, 0x12, 0x17, 0x0a, 0x13, 0x50,
	0x41, 0x52, 0x54, 0x49, 0x45, 0x53, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x41, 0x4c, 0x4c, 0x4f, 0x57,
	0x45, 0x44, 0x10, 0x05, 0x22, 0x35, 0x0a, 0x16, 0x44, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x42,
	0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b,
	0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x22, 0x19, 0x0a, 0x17, 0x44,
	0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xd0, 0x01, 0x0a, 0x1c, 0x44, 0x65, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5c, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x44, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d,
	0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x6b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x52, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x4f, 0x54, 0x5f, 0x49, 0x4e, 0x5f, 0x51,
	0x55, 0x45, 0x55, 0x45, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4e, 0x4f, 0x5f, 0x50, 0x45, 0x52,
	0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x41, 0x4c, 0x52,
	0x45, 0x41, 0x44, 0x59, 0x5f, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x44, 0x5f, 0x46, 0x4f, 0x52, 0x5f,
	0x44, 0x45, 0x51, 0x55, 0x45, 0x55, 0x45, 0x10, 0x02, 0x22, 0x50, 0x0a, 0x1a, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x56, 0x6f, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x61, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x61, 0x70, 0x49, 0x64, 0x22, 0x1d, 0x0a, 0x1b, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x56, 0x6f,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xb6, 0x01, 0x0a, 0x20, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x56, 0x6f,
	0x74, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x60, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x48, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x6b, 0x75,
	0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x56, 0x6f, 0x74, 0x65,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x22, 0x30, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x4d, 0x41, 0x50, 0x10,
	0x00, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x4f, 0x54, 0x5f, 0x49, 0x4e, 0x5f, 0x51, 0x55, 0x45, 0x55,
	0x45, 0x10, 0x01, 0x32, 0x8d, 0x03, 0x0a, 0x0a, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x6d, 0x61, 0x6b,
	0x65, 0x72, 0x12, 0x76, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x12, 0x30, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72, 0x74, 0x61,
	0x6c, 0x2e, 0x6b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x6b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x7c, 0x0a, 0x0f, 0x44, 0x65,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x32, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x6b, 0x75, 0x72, 0x75,
	0x73, 0x68, 0x69, 0x6d, 0x69, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x65, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x33, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e,
	0x6b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x44,
	0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x42, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x88, 0x01, 0x0a, 0x13, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x56, 0x6f, 0x74, 0x65,
	0x12, 0x36, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x6b,
	0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x70, 0x56, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x37, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x65,
	0x6d, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x6b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x4d, 0x61, 0x70, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x5c, 0x0a, 0x19, 0x64, 0x65, 0x76, 0x2e, 0x65, 0x6d, 0x6f, 0x72, 0x74,
	0x61, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69,
	0x42, 0x16, 0x4b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x46, 0x72, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6d, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x6d, 0x63,
	0x2f, 0x6b, 0x75, 0x72, 0x75, 0x73, 0x68, 0x69, 0x6d, 0x69, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_proto_rawDescOnce sync.Once
	file_grpc_proto_rawDescData = file_grpc_proto_rawDesc
)

func file_grpc_proto_rawDescGZIP() []byte {
	file_grpc_proto_rawDescOnce.Do(func() {
		file_grpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_proto_rawDescData)
	})
	return file_grpc_proto_rawDescData
}

var file_grpc_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_grpc_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_grpc_proto_goTypes = []interface{}{
	(QueueByPlayerErrorResponse_ErrorReason)(0),       // 0: dev.emortal.kurushimi.grpc.QueueByPlayerErrorResponse.ErrorReason
	(DequeueByPlayerErrorResponse_ErrorReason)(0),     // 1: dev.emortal.kurushimi.grpc.DequeueByPlayerErrorResponse.ErrorReason
	(ChangePlayerMapVoteErrorResponse_ErrorReason)(0), // 2: dev.emortal.kurushimi.grpc.ChangePlayerMapVoteErrorResponse.ErrorReason
	(*QueueByPlayerRequest)(nil),                      // 3: dev.emortal.kurushimi.grpc.QueueByPlayerRequest
	(*QueueByPlayerResponse)(nil),                     // 4: dev.emortal.kurushimi.grpc.QueueByPlayerResponse
	(*QueueByPlayerErrorResponse)(nil),                // 5: dev.emortal.kurushimi.grpc.QueueByPlayerErrorResponse
	(*DequeueByPlayerRequest)(nil),                    // 6: dev.emortal.kurushimi.grpc.DequeueByPlayerRequest
	(*DequeueByPlayerResponse)(nil),                   // 7: dev.emortal.kurushimi.grpc.DequeueByPlayerResponse
	(*DequeueByPlayerErrorResponse)(nil),              // 8: dev.emortal.kurushimi.grpc.DequeueByPlayerErrorResponse
	(*ChangePlayerMapVoteRequest)(nil),                // 9: dev.emortal.kurushimi.grpc.ChangePlayerMapVoteRequest
	(*ChangePlayerMapVoteResponse)(nil),               // 10: dev.emortal.kurushimi.grpc.ChangePlayerMapVoteResponse
	(*ChangePlayerMapVoteErrorResponse)(nil),          // 11: dev.emortal.kurushimi.grpc.ChangePlayerMapVoteErrorResponse
}
var file_grpc_proto_depIdxs = []int32{
	0,  // 0: dev.emortal.kurushimi.grpc.QueueByPlayerErrorResponse.reason:type_name -> dev.emortal.kurushimi.grpc.QueueByPlayerErrorResponse.ErrorReason
	1,  // 1: dev.emortal.kurushimi.grpc.DequeueByPlayerErrorResponse.reason:type_name -> dev.emortal.kurushimi.grpc.DequeueByPlayerErrorResponse.ErrorReason
	2,  // 2: dev.emortal.kurushimi.grpc.ChangePlayerMapVoteErrorResponse.reason:type_name -> dev.emortal.kurushimi.grpc.ChangePlayerMapVoteErrorResponse.ErrorReason
	3,  // 3: dev.emortal.kurushimi.grpc.Matchmaker.QueueByPlayer:input_type -> dev.emortal.kurushimi.grpc.QueueByPlayerRequest
	6,  // 4: dev.emortal.kurushimi.grpc.Matchmaker.DequeueByPlayer:input_type -> dev.emortal.kurushimi.grpc.DequeueByPlayerRequest
	9,  // 5: dev.emortal.kurushimi.grpc.Matchmaker.ChangePlayerMapVote:input_type -> dev.emortal.kurushimi.grpc.ChangePlayerMapVoteRequest
	4,  // 6: dev.emortal.kurushimi.grpc.Matchmaker.QueueByPlayer:output_type -> dev.emortal.kurushimi.grpc.QueueByPlayerResponse
	7,  // 7: dev.emortal.kurushimi.grpc.Matchmaker.DequeueByPlayer:output_type -> dev.emortal.kurushimi.grpc.DequeueByPlayerResponse
	10, // 8: dev.emortal.kurushimi.grpc.Matchmaker.ChangePlayerMapVote:output_type -> dev.emortal.kurushimi.grpc.ChangePlayerMapVoteResponse
	6,  // [6:9] is the sub-list for method output_type
	3,  // [3:6] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_grpc_proto_init() }
func file_grpc_proto_init() {
	if File_grpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueByPlayerRequest); i {
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
		file_grpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueByPlayerResponse); i {
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
		file_grpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueByPlayerErrorResponse); i {
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
		file_grpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DequeueByPlayerRequest); i {
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
		file_grpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DequeueByPlayerResponse); i {
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
		file_grpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DequeueByPlayerErrorResponse); i {
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
		file_grpc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangePlayerMapVoteRequest); i {
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
		file_grpc_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangePlayerMapVoteResponse); i {
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
		file_grpc_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangePlayerMapVoteErrorResponse); i {
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
	file_grpc_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_proto_goTypes,
		DependencyIndexes: file_grpc_proto_depIdxs,
		EnumInfos:         file_grpc_proto_enumTypes,
		MessageInfos:      file_grpc_proto_msgTypes,
	}.Build()
	File_grpc_proto = out.File
	file_grpc_proto_rawDesc = nil
	file_grpc_proto_goTypes = nil
	file_grpc_proto_depIdxs = nil
}
