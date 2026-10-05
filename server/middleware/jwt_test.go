package middleware

import (
	restful "github.com/emicklei/go-restful/v3"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testSecret = "test-secret-key-for-jwt"

func TestOptionalReadAuth(t *testing.T) {
	token, err := GenerateToken(42, "editor", testSecret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name, authorization string
		authenticated       bool
	}{
		{"anonymous reader", "", false}, {"invalid token remains public", "Bearer invalid", false}, {"authenticated editor", "Bearer " + token, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			container := restful.NewContainer()
			ws := new(restful.WebService).Path("/")
			ws.Route(ws.GET("/read").Filter(OptionalJWTFilter(testSecret)).To(func(req *restful.Request, resp *restful.Response) {
				authenticated := req.Attribute("claims") != nil
				if authenticated != test.authenticated {
					t.Errorf("reader authenticated = %v, want %v", authenticated, test.authenticated)
				}
				resp.WriteHeader(http.StatusNoContent)
			}))
			container.Add(ws)
			request := httptest.NewRequest(http.MethodGet, "/read", nil)
			request.Header.Set("Authorization", test.authorization)
			response := httptest.NewRecorder()
			container.ServeHTTP(response, request)
			if response.Code != http.StatusNoContent {
				t.Errorf("public read rejected: status %d", response.Code)
			}
		})
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	tokenStr, err := GenerateToken(42, "admin", testSecret, 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token string")
	}

	claims, err := ValidateToken(tokenStr, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("expected UserID 42, got %d", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("expected Username 'admin', got %s", claims.Username)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("expected non-nil ExpiresAt")
	}
	if claims.IssuedAt == nil {
		t.Fatal("expected non-nil IssuedAt")
	}
}

func TestValidateExpiredToken(t *testing.T) {
	tokenStr, err := GenerateToken(1, "user", testSecret, -1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ValidateToken(tokenStr, testSecret)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateInvalidToken(t *testing.T) {
	_, err := ValidateToken("this-is-not-a-valid-jwt", testSecret)
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}
