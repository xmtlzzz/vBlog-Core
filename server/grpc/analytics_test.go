package grpc

import (
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"vblog-core/model"
	pb "vblog-core/proto"
	"vblog-core/service"
	"vblog-core/testutil"
)

func startTestServer(t *testing.T) (*service.ChangeLogService, *service.SettingService, pb.BlogAnalyticsClient, func()) {
	t.Helper()
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

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	client := pb.NewBlogAnalyticsClient(conn)

	cleanup := func() {
		conn.Close()
		srv.GrpcServer.Stop()
	}
	return cl, st, client, cleanup
}

func TestGetLatestStats(t *testing.T) {
	_, _, client, cleanup := startTestServer(t)
	defer cleanup()

	stats, err := client.GetLatestStats(t.Context(), &pb.Empty{})
	if err != nil {
		t.Fatalf("GetLatestStats failed: %v", err)
	}
	if stats == nil {
		t.Fatal("expected non-nil stats")
	}
	// Stats fields should be zero-valued with empty DB (or valid counts)
	if stats.PvToday < 0 || stats.UvToday < 0 {
		t.Error("pv_today and uv_today should be non-negative")
	}
	if stats.TotalPosts < 0 || stats.TotalViews < 0 || stats.TotalComments < 0 || stats.TotalTags < 0 {
		t.Error("total counts should be non-negative")
	}
}

func TestGetTrends(t *testing.T) {
	_, _, client, cleanup := startTestServer(t)
	defer cleanup()

	resp, err := client.GetTrends(t.Context(), &pb.GetTrendsRequest{
		Granularity: "day",
		Count:       7,
	})
	if err != nil {
		t.Fatalf("GetTrends failed: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	// With empty daily_stats table, points may be empty
	if resp.Points == nil {
		resp.Points = []*pb.TrendPoint{} // normalize nil to empty
	}
}

func TestGetTrendsDefaultCount(t *testing.T) {
	_, _, client, cleanup := startTestServer(t)
	defer cleanup()

	// Count=0 should default to 7
	resp, err := client.GetTrends(t.Context(), &pb.GetTrendsRequest{
		Granularity: "day",
		Count:       0,
	})
	if err != nil {
		t.Fatalf("GetTrends with default count failed: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

func TestWatchChangesInvalidKey(t *testing.T) {
	_, _, client, cleanup := startTestServer(t)
	defer cleanup()

	stream, err := client.WatchChanges(t.Context(), &pb.WatchRequest{
		ApiKey:  "wrong-key",
		SinceId: 0,
	})
	if err != nil {
		t.Fatalf("WatchChanges call failed: %v", err)
	}

	// Should receive an error on first Recv due to invalid key
	_, err = stream.Recv()
	if err == nil {
		t.Fatal("expected error for invalid api key, got nil")
	}
}

func TestWatchChangesValidKey(t *testing.T) {
	cl, stSvc, client, cleanup := startTestServer(t)
	defer cleanup()

	// Set API key
	if err := stSvc.Set("grpc_api_key", "test-key-123"); err != nil {
		t.Fatalf("failed to set api key: %v", err)
	}

	// Get the latest change log ID so we only receive new events
	latestID, err := cl.GetLatestID()
	if err != nil {
		t.Fatalf("failed to get latest change log ID: %v", err)
	}

	stream, err := client.WatchChanges(t.Context(), &pb.WatchRequest{
		ApiKey:  "test-key-123",
		SinceId: latestID,
	})
	if err != nil {
		t.Fatalf("WatchChanges failed: %v", err)
	}

	// Write a change while watching
	if err := cl.Write("new_post", nil, "Stream Test Post", ""); err != nil {
		t.Fatalf("failed to write change log: %v", err)
	}

	// Receive event with timeout
	done := make(chan error, 1)
	var event *pb.ChangeEvent
	go func() {
		e, err := stream.Recv()
		event = e
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("stream.Recv failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for change event")
	}

	if event.Type != "new_post" {
		t.Errorf("expected 'new_post', got '%s'", event.Type)
	}

	// Cleanup test data
	db := testutil.GetTestDB(t)
	db.Where("title = ?", "Stream Test Post").Delete(&model.ChangeLog{})
	db.Where("key = ?", "grpc_api_key").Delete(&model.Setting{})
}
