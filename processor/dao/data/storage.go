package data

import (
	"io"
)

type DataStorage interface {
	Load(string) (io.ReadCloser, error)
	Close(io.ReadCloser) error
}
