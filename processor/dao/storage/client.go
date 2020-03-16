package storage

import (
	"context"
	"time"

	pb "github.com/kachan1208/something/api/proto/storage"
	"github.com/kachan1208/something/processor/model"
)

type Client struct {
	client pb.StorageClient
}

func NewClient(client pb.StorageClient) *Client {
	return &Client{
		client: client,
	}
}

func (s *Client) StoreBatchOfPorts(batch []*model.Port) (uint64, error) {
	count := uint64(len(batch))
	ports := make([]*pb.Port, 0, count)
	for i, port := range batch {
		ports[i] = &pb.Port{
			Key:         port.Code,
			Name:        port.Name,
			City:        port.City,
			Country:     port.Country,
			Alias:       port.Alias,
			Regions:     port.Regions,
			Coordinates: port.Coordinates,
			Province:    port.Province,
			Unlocs:      port.Unlocs,
			Code:        port.Code,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.client.StorePorts(
		ctx,
		&pb.StorePortsReq{Ports: ports},
	)

	return count, err
}
