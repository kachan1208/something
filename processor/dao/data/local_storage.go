package data

import (
	"io"
	"os"
	"path"
)

type LocalStorage struct {
	path string
}

func NewLocalStorage(path string) *LocalStorage {
	return &LocalStorage{
		path: path,
	}
}

func (l *LocalStorage) Load(name string) (io.ReadCloser, error) {
	file, err := os.Open(path.Join(l.path, name))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *LocalStorage) Close(file io.ReadCloser) error {
	return file.Close()
}
