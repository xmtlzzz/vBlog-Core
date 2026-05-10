package grpc

import (
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "vblog-core/proto"
)

func TestServer_StartAndPing(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := NewServer()
	go srv.GrpcServer.Serve(lis)
	defer srv.GrpcServer.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBlogAnalyticsClient(conn)
	_, err = client.Ping(t.Context(), &pb.Empty{})
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}
