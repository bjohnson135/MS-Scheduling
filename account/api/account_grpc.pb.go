// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: account.proto

package main

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AccountService_Create_FullMethodName                  = "/account.AccountService/Create"
	AccountService_List_FullMethodName                    = "/account.AccountService/List"
	AccountService_Get_FullMethodName                     = "/account.AccountService/Get"
	AccountService_Update_FullMethodName                  = "/account.AccountService/Update"
	AccountService_UpdatePassword_FullMethodName          = "/account.AccountService/UpdatePassword"
	AccountService_RequestPasswordReset_FullMethodName    = "/account.AccountService/RequestPasswordReset"
	AccountService_RequestEmailChange_FullMethodName      = "/account.AccountService/RequestEmailChange"
	AccountService_VerifyPassword_FullMethodName          = "/account.AccountService/VerifyPassword"
	AccountService_ChangeEmail_FullMethodName             = "/account.AccountService/ChangeEmail"
	AccountService_GetOrCreate_FullMethodName             = "/account.AccountService/GetOrCreate"
	AccountService_GetAccountByPhonenumber_FullMethodName = "/account.AccountService/GetAccountByPhonenumber"
	AccountService_TrackEvent_FullMethodName              = "/account.AccountService/TrackEvent"
	AccountService_SyncUser_FullMethodName                = "/account.AccountService/SyncUser"
)

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	Create(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*Account, error)
	List(ctx context.Context, in *GetAccountListRequest, opts ...grpc.CallOption) (*AccountList, error)
	Get(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*Account, error)
	Update(ctx context.Context, in *Account, opts ...grpc.CallOption) (*Account, error)
	UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RequestPasswordReset(ctx context.Context, in *PasswordResetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RequestEmailChange(ctx context.Context, in *EmailChangeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	VerifyPassword(ctx context.Context, in *VerifyPasswordRequest, opts ...grpc.CallOption) (*Account, error)
	ChangeEmail(ctx context.Context, in *EmailConfirmation, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetOrCreate(ctx context.Context, in *GetOrCreateRequest, opts ...grpc.CallOption) (*Account, error)
	GetAccountByPhonenumber(ctx context.Context, in *GetAccountByPhonenumberRequest, opts ...grpc.CallOption) (*Account, error)
	TrackEvent(ctx context.Context, in *TrackEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SyncUser(ctx context.Context, in *SyncUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) Create(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*Account, error) {
	out := new(Account)
	err := c.cc.Invoke(ctx, AccountService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) List(ctx context.Context, in *GetAccountListRequest, opts ...grpc.CallOption) (*AccountList, error) {
	out := new(AccountList)
	err := c.cc.Invoke(ctx, AccountService_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Get(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*Account, error) {
	out := new(Account)
	err := c.cc.Invoke(ctx, AccountService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Update(ctx context.Context, in *Account, opts ...grpc.CallOption) (*Account, error) {
	out := new(Account)
	err := c.cc.Invoke(ctx, AccountService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccountService_UpdatePassword_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) RequestPasswordReset(ctx context.Context, in *PasswordResetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccountService_RequestPasswordReset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) RequestEmailChange(ctx context.Context, in *EmailChangeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccountService_RequestEmailChange_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) VerifyPassword(ctx context.Context, in *VerifyPasswordRequest, opts ...grpc.CallOption) (*Account, error) {
	out := new(Account)
	err := c.cc.Invoke(ctx, AccountService_VerifyPassword_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) ChangeEmail(ctx context.Context, in *EmailConfirmation, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccountService_ChangeEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetOrCreate(ctx context.Context, in *GetOrCreateRequest, opts ...grpc.CallOption) (*Account, error) {
	out := new(Account)
	err := c.cc.Invoke(ctx, AccountService_GetOrCreate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetAccountByPhonenumber(ctx context.Context, in *GetAccountByPhonenumberRequest, opts ...grpc.CallOption) (*Account, error) {
	out := new(Account)
	err := c.cc.Invoke(ctx, AccountService_GetAccountByPhonenumber_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) TrackEvent(ctx context.Context, in *TrackEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccountService_TrackEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) SyncUser(ctx context.Context, in *SyncUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, AccountService_SyncUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations must embed UnimplementedAccountServiceServer
// for forward compatibility
type AccountServiceServer interface {
	Create(context.Context, *CreateAccountRequest) (*Account, error)
	List(context.Context, *GetAccountListRequest) (*AccountList, error)
	Get(context.Context, *GetAccountRequest) (*Account, error)
	Update(context.Context, *Account) (*Account, error)
	UpdatePassword(context.Context, *UpdatePasswordRequest) (*emptypb.Empty, error)
	RequestPasswordReset(context.Context, *PasswordResetRequest) (*emptypb.Empty, error)
	RequestEmailChange(context.Context, *EmailChangeRequest) (*emptypb.Empty, error)
	VerifyPassword(context.Context, *VerifyPasswordRequest) (*Account, error)
	ChangeEmail(context.Context, *EmailConfirmation) (*emptypb.Empty, error)
	GetOrCreate(context.Context, *GetOrCreateRequest) (*Account, error)
	GetAccountByPhonenumber(context.Context, *GetAccountByPhonenumberRequest) (*Account, error)
	TrackEvent(context.Context, *TrackEventRequest) (*emptypb.Empty, error)
	SyncUser(context.Context, *SyncUserRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedAccountServiceServer()
}

// UnimplementedAccountServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccountServiceServer struct {
}

func (UnimplementedAccountServiceServer) Create(context.Context, *CreateAccountRequest) (*Account, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedAccountServiceServer) List(context.Context, *GetAccountListRequest) (*AccountList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedAccountServiceServer) Get(context.Context, *GetAccountRequest) (*Account, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedAccountServiceServer) Update(context.Context, *Account) (*Account, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedAccountServiceServer) UpdatePassword(context.Context, *UpdatePasswordRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePassword not implemented")
}
func (UnimplementedAccountServiceServer) RequestPasswordReset(context.Context, *PasswordResetRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestPasswordReset not implemented")
}
func (UnimplementedAccountServiceServer) RequestEmailChange(context.Context, *EmailChangeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestEmailChange not implemented")
}
func (UnimplementedAccountServiceServer) VerifyPassword(context.Context, *VerifyPasswordRequest) (*Account, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyPassword not implemented")
}
func (UnimplementedAccountServiceServer) ChangeEmail(context.Context, *EmailConfirmation) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeEmail not implemented")
}
func (UnimplementedAccountServiceServer) GetOrCreate(context.Context, *GetOrCreateRequest) (*Account, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrCreate not implemented")
}
func (UnimplementedAccountServiceServer) GetAccountByPhonenumber(context.Context, *GetAccountByPhonenumberRequest) (*Account, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountByPhonenumber not implemented")
}
func (UnimplementedAccountServiceServer) TrackEvent(context.Context, *TrackEventRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TrackEvent not implemented")
}
func (UnimplementedAccountServiceServer) SyncUser(context.Context, *SyncUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncUser not implemented")
}
func (UnimplementedAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Create(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).List(ctx, req.(*GetAccountListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Get(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Account)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Update(ctx, req.(*Account))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_UpdatePassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).UpdatePassword(ctx, req.(*UpdatePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_RequestPasswordReset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).RequestPasswordReset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_RequestPasswordReset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).RequestPasswordReset(ctx, req.(*PasswordResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_RequestEmailChange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).RequestEmailChange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_RequestEmailChange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).RequestEmailChange(ctx, req.(*EmailChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_VerifyPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).VerifyPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_VerifyPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).VerifyPassword(ctx, req.(*VerifyPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_ChangeEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailConfirmation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).ChangeEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_ChangeEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).ChangeEmail(ctx, req.(*EmailConfirmation))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetOrCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetOrCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetOrCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetOrCreate(ctx, req.(*GetOrCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetAccountByPhonenumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountByPhonenumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetAccountByPhonenumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetAccountByPhonenumber_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetAccountByPhonenumber(ctx, req.(*GetAccountByPhonenumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_TrackEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).TrackEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_TrackEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).TrackEvent(ctx, req.(*TrackEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_SyncUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).SyncUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_SyncUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).SyncUser(ctx, req.(*SyncUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "account.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _AccountService_Create_Handler,
		},
		{
			MethodName: "List",
			Handler:    _AccountService_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _AccountService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _AccountService_Update_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _AccountService_UpdatePassword_Handler,
		},
		{
			MethodName: "RequestPasswordReset",
			Handler:    _AccountService_RequestPasswordReset_Handler,
		},
		{
			MethodName: "RequestEmailChange",
			Handler:    _AccountService_RequestEmailChange_Handler,
		},
		{
			MethodName: "VerifyPassword",
			Handler:    _AccountService_VerifyPassword_Handler,
		},
		{
			MethodName: "ChangeEmail",
			Handler:    _AccountService_ChangeEmail_Handler,
		},
		{
			MethodName: "GetOrCreate",
			Handler:    _AccountService_GetOrCreate_Handler,
		},
		{
			MethodName: "GetAccountByPhonenumber",
			Handler:    _AccountService_GetAccountByPhonenumber_Handler,
		},
		{
			MethodName: "TrackEvent",
			Handler:    _AccountService_TrackEvent_Handler,
		},
		{
			MethodName: "SyncUser",
			Handler:    _AccountService_SyncUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account.proto",
}
