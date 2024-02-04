package caching

import (
	"context"
	"errors"
	"time"
)

type Cacher interface {
	Write(ctx context.Context, options WriteOptions) error
	Read(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type ConnectionOptions struct {
	Host     string
	Password string
	Database int
}

type WriteOptions struct {
	Key   string
	Value interface{}
	TTL   time.Duration
}

var ErrResultIsNil = errors.New("result is nil")
