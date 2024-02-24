package gateway

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"distributed-file-storage/internal/domain/models"
	distributedStoragev1 "distributed-file-storage/protos/gen/go/distributedStorage"
)

type Gateway struct {
	log         *slog.Logger
	fUploader   FileUploader
	fGetter     FileGetter
	appProvider AppProvider
}

const (
	// ChunkSize is the size of the chunk.
	ChunkSize = 1 << 26
)

type ChunkInfo struct {
	Data     []byte
	Index    int
	Filename string
	Mimetype string
}

func (g *Gateway) Upload(stream distributedStoragev1.DistributedStorage_UploadServer) error {
	var wg sync.WaitGroup
	chunkCh := make(chan []byte)

	wg.Add(1)
	go func() {
		defer close(chunkCh)
		defer wg.Done()
		var buffer []byte
		var mimeType string
		var mimeTypeDetected = false
		var filename string
		var chunkIndex int

		for {
			data, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Println("failed to receive data:", err)
				return
			}

			if filename == "" {
				filename = data.GetFilename()
			}

			if !mimeTypeDetected && len(buffer) >= 512 {
				mimeType = http.DetectContentType(buffer[:512])
				mimeTypeDetected = true
			}

			buffer = append(buffer, data.Data...)
			if len(buffer) >= ChunkSize {
				chunkCh <- ChunkInfo{Data: buffer[:ChunkSize], Index: chunkIndex, Filename: filename, Mimetype: mimeType}
				buffer = buffer[ChunkSize:]
				chunkIndex++
			}
		}
	}()

	return nil
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
func (g *Gateway) UploadFile(ctx context.Context, filename, mimeType string, blob []byte) error { //nolint
	panic("Not implemented")
}

// GetFile downloads a file from the storage.
func (g *Gateway) GetFile(ctx context.Context, filename string) ([]byte, error) { //nolint
	panic("Not implemented")
}
