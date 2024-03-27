package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"distributed-file-storage/internal/domain/models"
	"distributed-file-storage/internal/services/backend"
	distributedStoragev1 "distributed-file-storage/protos/gen/go/distributedStorage"

	db "distributed-file-storage/internal/storage/postgres"
)

type Gateway struct {
	log         *slog.Logger
	fUploader   FileUploader
	fGetter     FileGetter
	appProvider AppProvider
	storage     *db.Storage
	backend     *backend.Backend
}

const (
	// ChunkSize is the size of the chunk.
	ChunkSize = 1 << 26
)

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
	storage *db.Storage,
	backendService *backend.Backend,
) *Gateway {
	return &Gateway{
		log:         log,
		fUploader:   fileUploader,
		fGetter:     fileGetter,
		appProvider: appProvider,
		storage:     storage,
		backend:     backendService,
	}
}

func (g *Gateway) Upload(stream distributedStoragev1.DistributedStorage_UploadServer) error {
	var chunkIndex uint
	var buffer []byte
	var mimeType string
	var mimeTypeDetected bool
	var filename string
	ctx := stream.Context()

	for {
		// Receive data from stream
		data, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// Handle the last chunk
			if len(buffer) > 0 {
				if uploadErr := g.UploadFile(ctx, filename, mimeType, chunkIndex, buffer); uploadErr != nil {
					return fmt.Errorf("error downloading last chunk %w", uploadErr)
				}
			}
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive data: %w", err)
		}

		// Initialize filename and mimeType on the first chunk
		if filename == "" {
			filename = data.GetFilename() // Assuming your message has a GetFilename method
			chunkIndex = 0
		}

		// Append received data to the buffer
		buffer = append(buffer, data.Data...)

		// Process full chunks
		for len(buffer) >= ChunkSize {
			currentChunk := buffer[:ChunkSize]
			buffer = buffer[ChunkSize:]

			if !mimeTypeDetected {
				mimeType = http.DetectContentType(currentChunk[:512])
				mimeTypeDetected = true
			}

			// Process the current chunk
			if err := g.UploadFile(ctx, filename, mimeType, chunkIndex, currentChunk); err != nil {
				return fmt.Errorf("error uploading chunk %d: %w", chunkIndex, err)
			}
			chunkIndex++
		}
	}

	return nil
}

// UploadFile uploads a file to the storage.
func (g *Gateway) UploadFile(ctx context.Context, filename, mimeType string, chunkIndex uint, blob []byte) error {
	var metaID int64
	var err error

	// TODO: must be rewritten. Not optimal, if something goes wrong we will have info in meta but dont have the blob in backend.
	if chunkIndex == 0 {
		metaID, err = g.storage.SaveMeta(ctx, filename, mimeType)
		if err != nil {
			return fmt.Errorf("failed to save metadata: %w", err)
		}
	}

	if err := g.backend.Upload(ctx, blob, chunkIndex, metaID); err != nil {
		return fmt.Errorf("failed to upload blob: %w", err)
	}

	return nil
}

// GetFile downloads a file from the storage.
func (g *Gateway) GetFile(ctx context.Context, filename string) ([]byte, error) {
	fmt.Println(ctx, filename)
	panic("Not implemented")
}
