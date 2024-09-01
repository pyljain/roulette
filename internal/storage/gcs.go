package storage

import (
	"context"
	"io"

	gstorage "cloud.google.com/go/storage"
)

type gcs struct {
	bkt *gstorage.BucketHandle
}

func NewGCS(bucket string) (*gcs, error) {
	ctx := context.Background()
	client, err := gstorage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bkt := client.Bucket(bucket)
	return &gcs{bkt: bkt}, nil

}

func (g *gcs) Read(ctx context.Context, fileName string) ([]byte, error) {
	obj := g.bkt.Object(fileName)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	fileContents, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return fileContents, nil

}
