package handler

import (
	"context"

	pb "github.com/kachan1208/something/api/proto/processor"
	"github.com/kachan1208/something/processor/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PortHandler struct {
	portProcessor *processor.PortProcessor
}

func NewPortHandler(p *processor.PortProcessor) *PortHandler {
	return &PortHandler{
		portProcessor: p,
	}
}

func (p *PortHandler) ProcessFile(ctx context.Context, req *pb.ProcessFileReq) (*pb.ProcessFileResp, error) {
	//req validation

	jobID, err := p.portProcessor.Process(req.Filename)
	if err != nil {
		return nil, status.New(codes.Unknown, err.Error()).Err()
	}

	return &pb.ProcessFileResp{
		JobId: jobID,
	}, nil
}
