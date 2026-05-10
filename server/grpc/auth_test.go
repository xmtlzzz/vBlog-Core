package grpc

import (
	"testing"

	"vblog-core/model"
	"vblog-core/service"
	"vblog-core/testutil"
)

func TestAuthInterceptor_ValidKey(t *testing.T) {
	db := testutil.GetTestDB(t)
	st := service.NewSettingService(db)
	st.Set("grpc_api_key", "valid-key")

	interceptor := NewAuthInterceptor(st)
	err := interceptor.ValidateKey("valid-key")
	if err != nil {
		t.Fatalf("expected valid key to pass: %v", err)
	}

	db.Where("key = ?", "grpc_api_key").Delete(&model.Setting{})
}

func TestAuthInterceptor_InvalidKey(t *testing.T) {
	db := testutil.GetTestDB(t)
	st := service.NewSettingService(db)
	st.Set("grpc_api_key", "valid-key")

	interceptor := NewAuthInterceptor(st)
	err := interceptor.ValidateKey("wrong-key")
	if err == nil {
		t.Fatal("expected invalid key to fail")
	}

	db.Where("key = ?", "grpc_api_key").Delete(&model.Setting{})
}

func TestAuthInterceptor_NoKey(t *testing.T) {
	db := testutil.GetTestDB(t)
	st := service.NewSettingService(db)

	interceptor := NewAuthInterceptor(st)
	err := interceptor.ValidateKey("anything")
	if err == nil {
		t.Fatal("expected no key configured to fail")
	}
}
