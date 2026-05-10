package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	pb "vblog-core/proto"
)

type Server struct {
	pb.UnimplementedBlogAnalyticsServer
	GrpcServer *grpc.Server
}

func NewServer() *Server {
	s := &Server{
		GrpcServer: grpc.NewServer(),
	}
	pb.RegisterBlogAnalyticsServer(s.GrpcServer, s)
	return s
}

func (s *Server) Start(lis net.Listener) error {
	return s.GrpcServer.Serve(lis)
}

func (s *Server) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
