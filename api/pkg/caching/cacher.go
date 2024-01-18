package caching

import "context"

type Cacher interface {
	Write(ctx context.Context, key string, value interface{}) error
	Read(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type ConnectionOptions struct {
	Host     string
	Password string
	Database int
}
