package api

import (
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"vblog-core/middleware"
	"vblog-core/service"
)

// AuthResource provides authentication endpoints.
type AuthResource struct {
	Service *service.AuthService
	Secret  string
	Auth    restful.FilterFunction
}

// Rate limits: login is stricter (brute force).
var (
	loginRateFilter = middleware.RateLimitFilter(time.Minute, 10)
)

// Register adds auth routes to the given WebService.
func (a *AuthResource) Register(ws *restful.WebService) {
	ws.Route(ws.POST("/api/auth/login").Filter(loginRateFilter).To(a.login).
		Doc("User login").
		Notes("Authenticates a user and returns JWT access and refresh tokens. Rate limited to 10 requests/min per IP.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"auth"}).
		Reads(loginRequest{}).
		Writes(TokenResponse{}).
		Returns(200, "OK", TokenResponse{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(429, "Too Many Requests", ErrorResponse{}))

	if a.Auth != nil {
		ws.Route(ws.GET("/api/auth/me").Filter(a.Auth).To(a.me).
			Doc("Get current user").
			Notes("Returns the authenticated user info from the token.").
			Metadata(restfulspec.KeyOpenAPITags, []string{"auth"}).
			Writes(UserInfoResponse{}).
			Returns(200, "OK", UserInfoResponse{}).
			Returns(401, "Unauthorized", ErrorResponse{}))
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AuthResource) login(req *restful.Request, resp *restful.Response) {
	var body loginRequest
	if err := req.ReadEntity(&body); err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, err := a.Service.Login(body.Username, body.Password)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	accessToken, err := middleware.GenerateToken(user.ID, user.Username, a.Secret, 24*time.Hour)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
		return
	}

	refreshToken, err := middleware.GenerateToken(user.ID, user.Username, a.Secret, 7*24*time.Hour)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, map[string]string{"error": "failed to generate refresh token"})
		return
	}

	resp.WriteEntity(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (a *AuthResource) me(req *restful.Request, resp *restful.Response) {
	claims, ok := req.Attribute("claims").(*middleware.Claims)
	if !ok || claims == nil {
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	resp.WriteEntity(UserInfoResponse{
		ID:       claims.UserID,
		Username: claims.Username,
	})
}
