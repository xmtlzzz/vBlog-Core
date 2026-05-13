package api

import "vblog-core/model"

// PostListResponse represents a paginated list of posts.
type PostListResponse struct {
	Data  []postResp `json:"data"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
}

// CommentListResponse represents a paginated list of comments.
type CommentListResponse struct {
	Data  []model.Comment `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
}

// TagListResponse represents a list of tags.
type TagListResponse struct {
	Data []model.Tag `json:"data"`
}

// ComponentListResponse represents a list of components.
type ComponentListResponse struct {
	Data []model.Component `json:"data"`
}

// MessageResponse represents a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// TokenResponse represents a JWT token response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// DashboardStatsResponse represents dashboard statistics.
type DashboardStatsResponse struct {
	TotalPosts    int64 `json:"total_posts"`
	TotalViews    int64 `json:"total_views"`
	TotalComments int64 `json:"total_comments"`
	TotalTags     int64 `json:"total_tags"`
}

// UploadResponse represents a file upload response.
type UploadResponse struct {
	URL string `json:"url"`
}
