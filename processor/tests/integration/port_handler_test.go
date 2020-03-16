//+build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/kachan1208/something/api/proto/processor"
	"github.com/kachan1208/something/processor/dao/data"
	"github.com/kachan1208/something/processor/handler"
	"github.com/kachan1208/something/processor/processor"
	"github.com/kachan1208/something/processor/stream/buffer"
	"github.com/kachan1208/something/processor/stream/decoder"
	"github.com/kachan1208/something/processor/tests/fixtures"
)

func initPortHandler() *handler.PortHandler {
	jsonDecoder := decoder.NewJSONStreamDecoder()
	interfaceBuffer := buffer.NewInterfaceStream(1024)
	localStorage := data.NewLocalStorage("../fixtures")
	portProcessor := processor.NewPortProcessor(jsonDecoder, localStorage, interfaceBuffer, &StorageClientMock{})
	portHandler := handler.NewPortHandler(portProcessor)

	return portHandler
}

func TestHandlerProcessFile(t *testing.T) {
	resp, err := initPortHandler().ProcessFile(
		context.Background(),
		&pb.ProcessFileReq{
			Filename: fixtures.TestFilename,
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.JobId)
}

func TestHandlerProcessFileNotFound(t *testing.T) {
	_, err := initPortHandler().ProcessFile(
		context.Background(),
		&pb.ProcessFileReq{
			Filename: "not_found",
		},
	)

	assert.Error(t, err)
	assert.Equal(t, codes.Unknown, status.Code(err))
}

//TODO: test for validation errors
