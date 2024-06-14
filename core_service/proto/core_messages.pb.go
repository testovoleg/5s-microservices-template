// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: core_messages.proto

package coreService

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InvoiceHandler struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID          string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Title       string `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=Description,proto3" json:"Description,omitempty"`
	Version     string `protobuf:"bytes,4,opt,name=Version,proto3" json:"Version,omitempty"`
}

func (x *InvoiceHandler) Reset() {
	*x = InvoiceHandler{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvoiceHandler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvoiceHandler) ProtoMessage() {}

func (x *InvoiceHandler) ProtoReflect() protoreflect.Message {
	mi := &file_core_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvoiceHandler.ProtoReflect.Descriptor instead.
func (*InvoiceHandler) Descriptor() ([]byte, []int) {
	return file_core_messages_proto_rawDescGZIP(), []int{0}
}

func (x *InvoiceHandler) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *InvoiceHandler) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *InvoiceHandler) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *InvoiceHandler) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type InvoiceHandlersListReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InvoiceHandlersListReq) Reset() {
	*x = InvoiceHandlersListReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvoiceHandlersListReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvoiceHandlersListReq) ProtoMessage() {}

func (x *InvoiceHandlersListReq) ProtoReflect() protoreflect.Message {
	mi := &file_core_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvoiceHandlersListReq.ProtoReflect.Descriptor instead.
func (*InvoiceHandlersListReq) Descriptor() ([]byte, []int) {
	return file_core_messages_proto_rawDescGZIP(), []int{1}
}

type InvoiceHandlersListRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Handlers []*InvoiceHandler `protobuf:"bytes,1,rep,name=Handlers,proto3" json:"Handlers,omitempty"`
}

func (x *InvoiceHandlersListRes) Reset() {
	*x = InvoiceHandlersListRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvoiceHandlersListRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvoiceHandlersListRes) ProtoMessage() {}

func (x *InvoiceHandlersListRes) ProtoReflect() protoreflect.Message {
	mi := &file_core_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvoiceHandlersListRes.ProtoReflect.Descriptor instead.
func (*InvoiceHandlersListRes) Descriptor() ([]byte, []int) {
	return file_core_messages_proto_rawDescGZIP(), []int{2}
}

func (x *InvoiceHandlersListRes) GetHandlers() []*InvoiceHandler {
	if x != nil {
		return x.Handlers
	}
	return nil
}

type Product struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductID   string                 `protobuf:"bytes,1,opt,name=ProductID,proto3" json:"ProductID,omitempty"`
	Name        string                 `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Description string                 `protobuf:"bytes,3,opt,name=Description,proto3" json:"Description,omitempty"`
	Price       float64                `protobuf:"fixed64,4,opt,name=Price,proto3" json:"Price,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
}

func (x *Product) Reset() {
	*x = Product{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Product) ProtoMessage() {}

func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_core_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Product.ProtoReflect.Descriptor instead.
func (*Product) Descriptor() ([]byte, []int) {
	return file_core_messages_proto_rawDescGZIP(), []int{3}
}

func (x *Product) GetProductID() string {
	if x != nil {
		return x.ProductID
	}
	return ""
}

func (x *Product) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Product) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Product) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Product) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Product) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type UpdateProductReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductID   string  `protobuf:"bytes,1,opt,name=ProductID,proto3" json:"ProductID,omitempty"`
	Name        string  `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Description string  `protobuf:"bytes,3,opt,name=Description,proto3" json:"Description,omitempty"`
	Price       float64 `protobuf:"fixed64,4,opt,name=Price,proto3" json:"Price,omitempty"`
}

func (x *UpdateProductReq) Reset() {
	*x = UpdateProductReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProductReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProductReq) ProtoMessage() {}

func (x *UpdateProductReq) ProtoReflect() protoreflect.Message {
	mi := &file_core_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProductReq.ProtoReflect.Descriptor instead.
func (*UpdateProductReq) Descriptor() ([]byte, []int) {
	return file_core_messages_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateProductReq) GetProductID() string {
	if x != nil {
		return x.ProductID
	}
	return ""
}

func (x *UpdateProductReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateProductReq) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateProductReq) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type UpdateProductRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductID string `protobuf:"bytes,1,opt,name=ProductID,proto3" json:"ProductID,omitempty"`
}

func (x *UpdateProductRes) Reset() {
	*x = UpdateProductRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProductRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProductRes) ProtoMessage() {}

func (x *UpdateProductRes) ProtoReflect() protoreflect.Message {
	mi := &file_core_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProductRes.ProtoReflect.Descriptor instead.
func (*UpdateProductRes) Descriptor() ([]byte, []int) {
	return file_core_messages_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateProductRes) GetProductID() string {
	if x != nil {
		return x.ProductID
	}
	return ""
}

type DeleteProductByIdReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductID string `protobuf:"bytes,1,opt,name=ProductID,proto3" json:"ProductID,omitempty"`
}

func (x *DeleteProductByIdReq) Reset() {
	*x = DeleteProductByIdReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteProductByIdReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteProductByIdReq) ProtoMessage() {}

func (x *DeleteProductByIdReq) ProtoReflect() protoreflect.Message {
	mi := &file_core_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteProductByIdReq.ProtoReflect.Descriptor instead.
func (*DeleteProductByIdReq) Descriptor() ([]byte, []int) {
	return file_core_messages_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteProductByIdReq) GetProductID() string {
	if x != nil {
		return x.ProductID
	}
	return ""
}

type DeleteProductByIdRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteProductByIdRes) Reset() {
	*x = DeleteProductByIdRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_messages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteProductByIdRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteProductByIdRes) ProtoMessage() {}

func (x *DeleteProductByIdRes) ProtoReflect() protoreflect.Message {
	mi := &file_core_messages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteProductByIdRes.ProtoReflect.Descriptor instead.
func (*DeleteProductByIdRes) Descriptor() ([]byte, []int) {
	return file_core_messages_proto_rawDescGZIP(), []int{7}
}

var File_core_messages_proto protoreflect.FileDescriptor

var file_core_messages_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x63, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x72, 0x0a, 0x0e, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x48, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a,
	0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x18, 0x0a, 0x16, 0x49, 0x6e, 0x76, 0x6f, 0x69,
	0x63, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x22, 0x51, 0x0a, 0x16, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x12, 0x37, 0x0a, 0x08, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x76, 0x6f,
	0x69, 0x63, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x52, 0x08, 0x48, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x73, 0x22, 0xe7, 0x01, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x44, 0x12, 0x12,
	0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x7c,
	0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52,
	0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x44,
	0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x22, 0x30, 0x0a, 0x10,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x44, 0x22, 0x34,
	0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x42,
	0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x49, 0x44, 0x22, 0x16, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x42, 0x10, 0x5a, 0x0e,
	0x2e, 0x2f, 0x3b, 0x63, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_messages_proto_rawDescOnce sync.Once
	file_core_messages_proto_rawDescData = file_core_messages_proto_rawDesc
)

func file_core_messages_proto_rawDescGZIP() []byte {
	file_core_messages_proto_rawDescOnce.Do(func() {
		file_core_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_messages_proto_rawDescData)
	})
	return file_core_messages_proto_rawDescData
}

var file_core_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_core_messages_proto_goTypes = []interface{}{
	(*InvoiceHandler)(nil),         // 0: coreService.InvoiceHandler
	(*InvoiceHandlersListReq)(nil), // 1: coreService.InvoiceHandlersListReq
	(*InvoiceHandlersListRes)(nil), // 2: coreService.InvoiceHandlersListRes
	(*Product)(nil),                // 3: coreService.Product
	(*UpdateProductReq)(nil),       // 4: coreService.UpdateProductReq
	(*UpdateProductRes)(nil),       // 5: coreService.UpdateProductRes
	(*DeleteProductByIdReq)(nil),   // 6: coreService.DeleteProductByIdReq
	(*DeleteProductByIdRes)(nil),   // 7: coreService.DeleteProductByIdRes
	(*timestamppb.Timestamp)(nil),  // 8: google.protobuf.Timestamp
}
var file_core_messages_proto_depIdxs = []int32{
	0, // 0: coreService.InvoiceHandlersListRes.Handlers:type_name -> coreService.InvoiceHandler
	8, // 1: coreService.Product.CreatedAt:type_name -> google.protobuf.Timestamp
	8, // 2: coreService.Product.UpdatedAt:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_core_messages_proto_init() }
func file_core_messages_proto_init() {
	if File_core_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_core_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvoiceHandler); i {
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
		file_core_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvoiceHandlersListReq); i {
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
		file_core_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvoiceHandlersListRes); i {
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
		file_core_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Product); i {
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
		file_core_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProductReq); i {
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
		file_core_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProductRes); i {
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
		file_core_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteProductByIdReq); i {
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
		file_core_messages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteProductByIdRes); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_messages_proto_goTypes,
		DependencyIndexes: file_core_messages_proto_depIdxs,
		MessageInfos:      file_core_messages_proto_msgTypes,
	}.Build()
	File_core_messages_proto = out.File
	file_core_messages_proto_rawDesc = nil
	file_core_messages_proto_goTypes = nil
	file_core_messages_proto_depIdxs = nil
}
