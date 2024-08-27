package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManagerConf struct {
	Secret []byte
}

type JWTManager struct {
	conf *JWTManagerConf
}

func NewJWTManager(conf *JWTManagerConf) *JWTManager {
	return &JWTManager{
		conf: conf,
	}
}

type ClaimsPayload struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	RoleName string `json:"user_role"`
}

type Claims struct {
	Payload *ClaimsPayload
	jwt.RegisteredClaims
}

func (j *JWTManager) GenerateToken(paload *ClaimsPayload) (string, error) {
	claims := Claims{
		Payload: paload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.conf.Secret)
}

func (j *JWTManager) ParseToken(tokenString string) (*ClaimsPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.conf.Secret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Payload, nil
	} else {
		return nil, err
	}
}
