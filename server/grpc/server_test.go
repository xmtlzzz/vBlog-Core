package grpc

import (
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "vblog-core/proto"
	"vblog-core/service"
	"vblog-core/testutil"
)

func TestServer_StartAndPing(t *testing.T) {
	db := testutil.GetTestDB(t)
	ds := service.NewDailyStatsService(db)
	cl := service.NewChangeLogService(db)
	pv := service.NewPageViewService(db)
	st := service.NewSettingService(db)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := NewServer(ds, cl, pv, st)
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
