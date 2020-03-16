package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthPb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"

	pbProc "github.com/kachan1208/something/api/proto/processor"
	pbStor "github.com/kachan1208/something/api/proto/storage"
	"github.com/kachan1208/something/processor/config"
	"github.com/kachan1208/something/processor/dao/data"
	"github.com/kachan1208/something/processor/dao/storage"
	"github.com/kachan1208/something/processor/handler"
	"github.com/kachan1208/something/processor/processor"
	"github.com/kachan1208/something/processor/stream/buffer"
	"github.com/kachan1208/something/processor/stream/decoder"
)

//di can be useful here
//monitoring(readiness endpoint with dependencies monitoring)
//metrics(rpc, latency, statuses...)
//proper logging

func main() {
	decoder := decoder.NewJSONStreamDecoder()
	localStorage := data.NewLocalStorage("/")
	buf := buffer.NewInterfaceStream(1024)

	conn, _ := grpc.Dial(
		config.StorageAddress,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*5),
	)
	// if err != nil {
	// 	log.Fatal(conn)
	// }

	storageClient := storage.NewClient(pbStor.NewStorageClient(conn))
	portProcessor := processor.NewPortProcessor(decoder, localStorage, buf, storageClient)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	portsHandler := handler.NewPortHandler(portProcessor)
	portsGRPCServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime: 10 * time.Second,
			},
		),
	)
	pbProc.RegisterProcessorServer(portsGRPCServer, portsHandler)
	g.Go(func() error {
		listener, err := net.Listen("tcp", config.ProcessorAddress)
		if err != nil {
			return err
		}
		return portsGRPCServer.Serve(listener)
	})

	healthSrv := health.NewServer()
	healthGRPCServer := grpc.NewServer()
	healthPb.RegisterHealthServer(healthGRPCServer, healthSrv)
	g.Go(func() error {
		listener, err := net.Listen("tcp", config.ProcessorHealthAddress)
		if err != nil {
			return err
		}

		return healthGRPCServer.Serve(listener)
	})

	select {
	case <-signals:
		break
	case <-ctx.Done():
		break
	}

	fmt.Println("Graceful shutdown")

	cancel()

	healthSrv.SetServingStatus("service", healthPb.HealthCheckResponse_NOT_SERVING)
	portsGRPCServer.GracefulStop()
	healthGRPCServer.GracefulStop()
}
