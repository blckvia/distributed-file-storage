package app

import (
	"log/slog"

	grpcapp "distributed-file-storage/internal/app/grpc"
	"distributed-file-storage/internal/services/gateway"
	"distributed-file-storage/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	// TODO: initialize storage

	// TODO: init gateway service (gateway)
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}
	gatewayService := gateway.New(log, storage, storage, storage) //nolint
	//
	grpcApp := grpcapp.New(log, gatewayService, grpcPort) // nolint

	return &App{
		GRPCSrv: grpcApp,
	}
}
