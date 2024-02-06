package gateway

import (
	"context"
	"log/slog"

	"distributed-file-storage/internal/domain/models"
)

type Gateway struct {
	log         *slog.Logger
	fUploader   FileUploader
	fGetter     FileGetter
	appProvider AppProvider
}

type FileUploader interface {
	UploadFile(ctx context.Context, filename string, mimeType string, blob []byte) error
}

type FileGetter interface {
	GetFile(ctx context.Context, filename string) ([]byte, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New returns a new instance of the Gateway service.
func New(
	log *slog.Logger,
	fileUploader FileUploader,
	fileGetter FileGetter,
	appProvider AppProvider,
) *Gateway {
	return &Gateway{
		log:         log,
		fUploader:   fileUploader,
		fGetter:     fileGetter,
		appProvider: appProvider,
	}
}

// UploadFile uploads a file to the storage.
func (g *Gateway) UploadFile(ctx context.Context, filename, mime_type string, blob []byte) error {
	panic("Not implemented")
}

// GetFile downloads a file from the storage.
func (g *Gateway) GetFile(ctx context.Context, filename string) ([]byte, error) {
	panic("Not implemented")
}
