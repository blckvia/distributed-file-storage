package gateway

import (
	"context"
	"distributed-file-storage/internal/domain/models"
	distributedStoragev1 "distributed-file-storage/protos/gen/go/distributedStorage"
	"errors"
	"github.com/docker/docker/api/server/router/grpc"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
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

// 10mb ->

func (g *Gateway) Upload(stream distributedStoragev1.DistributedStorage_UploadServer) error {
	chunkCh := make(chan []byte)

	var buffer []byte
	var mimeType string
	var mimeTypeDetected = false
	var filename string
	var chunkIndex int

	for {
		var partData []byte
		chunkIndex = 0
		var tmp []byte
		if len(partData) < ChunkSize {
			data, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Println("failed to receive data:", err)
				break //
			}
			tmp = data.Data
			tmp = partData + data.Data[:ChunkSize - len([partData])]
		} else {
			tmp = partData
		}
		ctx := stream.Context()

		if filename == "" {
			filename = data.GetFilename()
		}

		buffer = tmp[:ChunkSize]
		if len(tmp) > ChunkSize {
			// тут подумать тут правильно разделить
			partData = tmp[ChunkSize+1:] // 0x11, 0x13, 0x13 0 0 0 0 0 0
		} else if len(tmp) <= ChunkSize {
			partData = make([]byte, 0, ChunkSize)
		}

		if chunkIndex == 0 && len(buffer) >= 512 {
			mimeType = http.DetectContentType(buffer[:512])
		}

		if err := g.fUploader.UploadFile(ctx, filename, mimeType, buffer); err != nil {
			return err
		}
		chunkIndex++
	}

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
func (g *Gateway) UploadFile(ctx context.Context, filename, mimeType string, chunkindex uint, blob []byte) error { //nolint
	if chunkindex == 0 {
		db.SaveMeta(ctx, filename, mimtype)
	}
	backend := grpc.Backend()
	backend.Upload(ctx, blob)
	db.SaveMeta(ctx, chunkindex, backend.ID)
}

// GetFile downloads a file from the storage.
func (g *Gateway) GetFile(ctx context.Context, filename string) ([]byte, error) { //nolint
	panic("Not implemented")
}


// UploadFile gtpc file backend
func (g *Gateway) UploadBlob(ctx context.Context, grcpClient) error { //nolint
	data := call.GetData()
	os.WriteFile("tmp-chunck-ID=jncdsnkdcns", data)
}
