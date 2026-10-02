package service

import (
	"errors"
	"time"

	"erp_monolith/backend/internal/modules/system/user/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtCustomClaims struct {
	UserID    uuid.UUID  `json:"uid"`
	CompanyID uuid.UUID  `json:"cid"`
	BranchID  *uuid.UUID `json:"bid,omitempty"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	jwt.RegisteredClaims
}

type jwtTokenService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTTokenService creates a TokenService implementation powered by golang-jwt/v5
func NewJWTTokenService(secret string, tokenDuration time.Duration) domain.TokenService {
	if tokenDuration <= 0 {
		tokenDuration = 24 * time.Hour
	}
	return &jwtTokenService{
		secretKey:     []byte(secret),
		tokenDuration: tokenDuration,
	}
}

func (s *jwtTokenService) GenerateToken(user *domain.User) (string, time.Time, error) {
	if user == nil {
		return "", time.Time{}, errors.New("user cannot be nil")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(s.tokenDuration)

	claims := jwtCustomClaims{
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		BranchID:  user.BranchID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *jwtTokenService) ValidateToken(tokenStr string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return &domain.Claims{
		UserID:    claims.UserID,
		CompanyID: claims.CompanyID,
		BranchID:  claims.BranchID,
		Email:     claims.Email,
		Name:      claims.Name,
		Role:      claims.Role,
	}, nil
}
