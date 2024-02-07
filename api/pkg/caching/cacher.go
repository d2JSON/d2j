package caching

import (
	"context"
	"errors"
	"time"
)

// Cacher represents an interface for working with cache.
type Cacher interface {
	Write(ctx context.Context, options WriteOptions) error
	Read(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

// ConnectionOptions represents an options that used for connecting to cache.
type ConnectionOptions struct {
	Host     string
	Password string
	Database int
}

// WriteOptions represents an options that used for writeing data into cache.
type WriteOptions struct {
	Key   string
	Value interface{}
	TTL   time.Duration
}

// ErrResultIsNil happens when cache returns nothing on read request.
var ErrResultIsNil = errors.New("result is nil")
