package models

type Meta struct {
	ID                 int64
	Filename           string
	MimeType           string
	BlobSequenceNumber []byte
}
