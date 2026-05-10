package main

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "vblog-core/proto"
)

type App struct {
	conn   *grpc.ClientConn
	client pb.BlogAnalyticsClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{ctx: ctx, cancel: cancel}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Connect(addr, apiKey string) error {
	if addr == "" {
		return errors.New("address is required")
	}
	// Probe address reachability since grpc.NewClient is lazy.
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	grpcConn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	a.conn = grpcConn
	a.client = pb.NewBlogAnalyticsClient(grpcConn)
	return nil
}

func (a *App) Disconnect() error {
	a.cancel()
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *App) GetLatestStats() (*pb.LatestStats, error) {
	if a.client == nil {
		return nil, errors.New("not connected")
	}
	return a.client.GetLatestStats(a.ctx, &pb.Empty{})
}

func (a *App) GetTrends(granularity string, count int32) (*pb.GetTrendsResponse, error) {
	if a.client == nil {
		return nil, errors.New("not connected")
	}
	return a.client.GetTrends(a.ctx, &pb.GetTrendsRequest{
		Granularity: granularity,
		Count:       count,
	})
}

func (a *App) WatchChanges(apiKey string, sinceID int64) error {
	if a.client == nil {
		return errors.New("not connected")
	}

	stream, err := a.client.WatchChanges(a.ctx, &pb.WatchRequest{
		ApiKey:  apiKey,
		SinceId: sinceID,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			event, err := stream.Recv()
			if err != nil {
				return
			}
			runtime.EventsEmit(a.ctx, "change", event)
		}
	}()
	return nil
}
