package grpc

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"vblog-core/model"
	pb "vblog-core/proto"
	"vblog-core/service"
)

type Server struct {
	pb.UnimplementedBlogAnalyticsServer
	GrpcServer    *grpc.Server
	DailyStatsSvc *service.DailyStatsService
	ChangeLogSvc  *service.ChangeLogService
	PageViewSvc   *service.PageViewService
	SettingSvc    *service.SettingService
}

func NewServer(ds *service.DailyStatsService, cl *service.ChangeLogService, pv *service.PageViewService, st *service.SettingService) *Server {
	s := &Server{
		GrpcServer:    grpc.NewServer(),
		DailyStatsSvc: ds,
		ChangeLogSvc:  cl,
		PageViewSvc:   pv,
		SettingSvc:    st,
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

func (s *Server) GetLatestStats(ctx context.Context, in *pb.Empty) (*pb.LatestStats, error) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	pvToday, uvToday, _ := s.PageViewSvc.GetPVUVByDate(today)
	pvYesterday, uvYesterday, _ := s.PageViewSvc.GetPVUVByDate(yesterday)

	var postCount, viewTotal, commentCount, tagCount int64
	s.DailyStatsSvc.DB.Model(&model.Post{}).Where("status = ?", "published").Count(&postCount)
	s.DailyStatsSvc.DB.Model(&model.Post{}).Select("COALESCE(SUM(views), 0)").Scan(&viewTotal)
	s.DailyStatsSvc.DB.Model(&model.Comment{}).Count(&commentCount)
	s.DailyStatsSvc.DB.Model(&model.Tag{}).Count(&tagCount)

	return &pb.LatestStats{
		PvToday:       pvToday,
		UvToday:       uvToday,
		TotalPosts:    postCount,
		TotalViews:    viewTotal,
		TotalComments: commentCount,
		TotalTags:     tagCount,
		PvYesterday:   pvYesterday,
		UvYesterday:   uvYesterday,
	}, nil
}

func (s *Server) GetTrends(ctx context.Context, in *pb.GetTrendsRequest) (*pb.GetTrendsResponse, error) {
	count := int(in.Count)
	if count <= 0 {
		count = 7
	}

	stats, err := s.DailyStatsSvc.GetTrends(in.Granularity, count)
	if err != nil {
		return nil, err
	}

	points := make([]*pb.TrendPoint, len(stats))
	for i, st := range stats {
		points[i] = &pb.TrendPoint{
			Label:        st.StatDate.Format("2006-01-02"),
			Pv:           st.PV,
			Uv:           st.UV,
			ViewTotal:    st.ViewTotal,
			CommentCount: int64(st.CommentCount),
			PostCount:    int64(st.PostCount),
		}
	}

	return &pb.GetTrendsResponse{Points: points}, nil
}

func (s *Server) WatchChanges(in *pb.WatchRequest, stream pb.BlogAnalytics_WatchChangesServer) error {
	// Validate API key
	apiKey, _ := s.SettingSvc.Get("grpc_api_key")
	if apiKey == "" || apiKey != in.ApiKey {
		return status.Error(codes.Unauthenticated, "invalid api_key")
	}

	// Catch-up: send missed changes
	sinceID := in.SinceId
	logs, err := s.ChangeLogSvc.GetAfterID(sinceID)
	if err != nil {
		return status.Error(codes.Internal, "failed to read change log")
	}
	for _, l := range logs {
		event := &pb.ChangeEvent{
			Id:        int64(l.ID),
			Type:      l.ChangeType,
			Title:     l.Title,
			Detail:    l.Detail,
			Timestamp: l.CreatedAt.Format(time.RFC3339),
		}
		if err := stream.Send(event); err != nil {
			return err
		}
		sinceID = int64(l.ID)
	}

	// Watch loop
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			newLogs, err := s.ChangeLogSvc.GetAfterID(sinceID)
			if err != nil {
				log.Printf("WatchChanges: error fetching logs after %d: %v", sinceID, err)
				continue
			}
			for _, l := range newLogs {
				event := &pb.ChangeEvent{
					Id:        int64(l.ID),
					Type:      l.ChangeType,
					Title:     l.Title,
					Detail:    l.Detail,
					Timestamp: l.CreatedAt.Format(time.RFC3339),
				}
				if err := stream.Send(event); err != nil {
					return err
				}
				sinceID = int64(l.ID)
			}
		}
	}
}
