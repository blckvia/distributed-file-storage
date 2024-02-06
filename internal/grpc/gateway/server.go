package gateway

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	distributedStoragev1 "distributed-file-storage/protos/gen/go/distributedStorage"
)

type Gateway interface {
	GetFile(ctx context.Context, filename string) ([]byte, error)
	Upload(stream distributedStoragev1.DistributedStorage_UploadServer) error
	mustEmbedUnimplementedDistributedStorageServer()
}

type ServerAPI struct {
	distributedStoragev1.UnimplementedDistributedStorageServer
	gateway Gateway
}

func Register(gRPC *grpc.Server, gateway Gateway) {
	distributedStoragev1.RegisterDistributedStorageServer(gRPC, &ServerAPI{gateway: gateway})
}

func (s *ServerAPI) GetFile(ctx context.Context, req *distributedStoragev1.GetfileRequest) (*distributedStoragev1.GetfileResponse, error) {
	if req.GetFilename() == "" {
		return nil, status.Error(codes.InvalidArgument, "filename is empty")
	}

	// TODO: implement file getter via gateway service
	data, err := s.gateway.GetFile(ctx, req.GetFilename())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &distributedStoragev1.GetfileResponse{
		Data: data,
	}, nil
}

func (s *ServerAPI) Upload(stream distributedStoragev1.DistributedStorage_UploadServer) error {
	panic("Not implemented")
}

func (s *ServerAPI) mustEmbedUnimplementedDistributedStorageServer() {
	panic("Not implemented")
}
