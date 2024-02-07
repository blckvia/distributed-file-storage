package app

import (
	"log/slog"

	grpcapp "distributed-file-storage/internal/app/grpc"
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

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
