package storage

import (
	"github.com/kachan1208/something/processor/model"
)

type StorageClient interface {
	StoreBatchOfPorts([]*model.Port) (uint64, error)
}
