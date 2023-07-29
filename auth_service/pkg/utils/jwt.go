package utils

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/paw1a/grpc-media-converter/auth_service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 5)
	return string(bytes)
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type JwtConfig struct {
	SecretKey       string `yaml:"secretKey"`
	ExpirationHours int64  `yaml:"expirationHours"`
}

type JwtProvider struct {
	cfg *JwtConfig
}

type JwtClaims struct {
	jwt.StandardClaims
	Id    int64
	Email string
}

func NewJwtProvider(cfg *JwtConfig) *JwtProvider {
	return &JwtProvider{cfg: cfg}
}

func (j *JwtProvider) GenerateToken(user domain.User) (signedToken string, err error) {
	claims := &JwtClaims{
		Id:    user.Id,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.cfg.ExpirationHours)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(j.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JwtProvider) ValidateToken(signedToken string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.cfg.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("jwt token is expired")
	}

	return claims, nil
}
