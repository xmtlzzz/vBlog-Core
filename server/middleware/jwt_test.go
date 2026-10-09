package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful/v3"
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

func TestJWTFilter_ExpiredAndMissing(t *testing.T) {
	container := restful.NewContainer()
	ws := new(restful.WebService).Path("/").Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/protected").Filter(JWTFilter(testSecret)).To(func(req *restful.Request, resp *restful.Response) {
		resp.WriteHeader(http.StatusOK)
	}))
	container.Add(ws)

	// 1. Missing token
	req1 := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp1 := httptest.NewRecorder()
	container.ServeHTTP(resp1, req1)
	if resp1.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing token, got %d", resp1.Code)
	}

	// 2. Expired token
	expToken, _ := GenerateToken(1, "testuser", testSecret, -1*time.Hour)
	req2 := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req2.Header.Set("Authorization", "Bearer "+expToken)
	resp2 := httptest.NewRecorder()
	container.ServeHTTP(resp2, req2)
	if resp2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", resp2.Code)
	}
	if !strings.Contains(resp2.Body.String(), "token expired") {
		t.Errorf("expected body to contain 'token expired', got %s", resp2.Body.String())
	}
}
