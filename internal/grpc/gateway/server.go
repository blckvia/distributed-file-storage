package gateway

import (
	"context"
	distributedStoragev1 "distributed-file-storage/protos/gen/go/distributedStorage"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	distributedStoragev1.UnimplementedDistributedStorageServer
}

func Register(gRPC *grpc.Server) {
	distributedStoragev1.RegisterDistributedStorageServer(gRPC, &ServerAPI{})
}

func (s *ServerAPI) GetFile(ctx context.Context, req *distributedStoragev1.GetfileRequest) (*distributedStoragev1.GetfileResponse, error) {
	panic("Not implemented")
}

func (s *ServerAPI) Upload(stream distributedStoragev1.DistributedStorage_UploadServer) error {
	panic("Not implemented")
}

func (s *ServerAPI) mustEmbedUnimplementedDistributedStorageServer() {
	panic("Not implemented")
}
