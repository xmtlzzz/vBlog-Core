package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt/v5"
)

// Claims extends jwt.RegisteredClaims with user identity fields.
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token for the given user.
func GenerateToken(userID uint, username, secret string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken parses and validates a JWT token string, returning its claims.
func ValidateToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTFilter returns a go-restful filter that validates Bearer tokens.
func JWTFilter(secret string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		auth := req.HeaderParameter("Authorization")
		if auth == "" {
			resp.WriteHeaderAndEntity(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := ValidateToken(tokenStr, secret)
		if err != nil {
			resp.WriteHeaderAndEntity(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}
		req.SetAttribute("claims", claims)
		chain.ProcessFilter(req, resp)
	}
}
