package storage

import "context"

type Storage interface {
	Read(ctx context.Context, fileName string) ([]byte, error)
}
