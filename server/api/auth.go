package api

import (
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"vblog-core/middleware"
	"vblog-core/service"
)

// AuthResource provides authentication endpoints.
type AuthResource struct {
	Service *service.AuthService
	Secret  string
}

// Register adds auth routes to the given WebService.
func (a *AuthResource) Register(ws *restful.WebService) {
	ws.Route(ws.POST("/api/auth/login").To(a.login).
		Doc("user login"))
	ws.Route(ws.POST("/api/auth/register").To(a.register).
		Doc("register new user"))
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

func (a *AuthResource) register(req *restful.Request, resp *restful.Response) {
	var body registerRequest
	if err := req.ReadEntity(&body); err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, err := a.Service.Register(body.Username, body.Password, body.Email)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, map[string]string{"error": "注册失败，用户名可能已存在"})
		return
	}

	accessToken, _ := middleware.GenerateToken(user.ID, user.Username, a.Secret, 24*time.Hour)
	refreshToken, _ := middleware.GenerateToken(user.ID, user.Username, a.Secret, 7*24*time.Hour)

	resp.WriteHeaderAndEntity(http.StatusCreated, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
