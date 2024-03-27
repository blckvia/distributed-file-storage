package backend

import (
	"context"
	"fmt"
	"log/slog"
)

type Backend struct {
	log       *slog.Logger
	bUploader BlobUploader
	bGetter   BlobGetter
}

type BlobUploader interface {
	Upload(ctx context.Context, blob []byte) error
}

type BlobGetter interface {
	Get(ctx context.Context, blobID int64) ([]byte, error)
}

// New returns a new instance of the Backend service.
func New(
	log *slog.Logger,
	blobUploader BlobUploader,
	blobGetter BlobGetter,
) *Backend {
	return &Backend{
		log:       log,
		bUploader: blobUploader,
		bGetter:   blobGetter,
	}
}

// Upload uploads a blob to the backend.
func (b *Backend) Upload(ctx context.Context, blob []byte, index uint, metaID int64) error {
	fmt.Println(index, metaID)
	return b.bUploader.Upload(ctx, blob)
}

// Get downloads a blob from the backend.
func (b *Backend) Get(ctx context.Context, blobID int64) ([]byte, error) {
	return b.bGetter.Get(ctx, blobID)
}
