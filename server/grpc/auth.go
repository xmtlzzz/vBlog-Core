package grpc

import (
	"errors"

	"vblog-core/service"
)

type AuthInterceptor struct {
	SettingSvc *service.SettingService
}

func NewAuthInterceptor(st *service.SettingService) *AuthInterceptor {
	return &AuthInterceptor{SettingSvc: st}
}

func (a *AuthInterceptor) ValidateKey(apiKey string) error {
	stored, err := a.SettingSvc.Get("grpc_api_key")
	if err != nil || stored == "" {
		return errors.New("grpc api key not configured")
	}
	if stored != apiKey {
		return errors.New("invalid api key")
	}
	return nil
}
