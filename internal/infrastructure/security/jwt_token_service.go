package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"auth-go/internal/domain/entity"
	"auth-go/internal/domain/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTTokenService implements TokenService using JWT
type JWTTokenService struct {
	secretKey          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	issuer             string
}

// Claims represents custom JWT claims
type Claims struct {
	UserID uuid.UUID     `json:"user_id"`
	Email  string        `json:"email"`
	Roles  []entity.Role `json:"roles"`
	jwt.RegisteredClaims
}

// NewJWTTokenService creates a new JWT token service
func NewJWTTokenService(
	secretKey string,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
	issuer string,
) *JWTTokenService {
	return &JWTTokenService{
		secretKey:          []byte(secretKey),
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
		issuer:             issuer,
	}
}

// GenerateAccessToken generates a JWT access token
func (s *JWTTokenService) GenerateAccessToken(claims service.TokenClaims) (string, error) {
	now := time.Now()
	jwtClaims := Claims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Roles:  claims.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpiry)),
			Issuer:    s.issuer,
			Subject:   claims.UserID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(s.secretKey)
}

// GenerateRefreshToken generates a cryptographically secure refresh token
func (s *JWTTokenService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ValidateAccessToken validates and parses an access token
func (s *JWTTokenService) ValidateAccessToken(tokenString string) (*service.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &service.TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Roles:  claims.Roles,
	}, nil
}

// GetAccessTokenExpiry returns the access token expiry duration
func (s *JWTTokenService) GetAccessTokenExpiry() time.Duration {
	return s.accessTokenExpiry
}

// GetRefreshTokenExpiry returns the refresh token expiry duration
func (s *JWTTokenService) GetRefreshTokenExpiry() time.Duration {
	return s.refreshTokenExpiry
}
