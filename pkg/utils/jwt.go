package utils

import (
	"fmt"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaim struct {
	UserID         string
	DeviceID       string
	PushNotifToken string
	PurchaseID     string
	Timezone       string
	PhoneLanguage  string
	jwt.RegisteredClaims
}

type JwtWrapper struct {
	SecretKey string
	Issuer    string
	Expire    int
}

type JwtClaim struct {
	ID       string
	UserName string
	UserId   string
	jwt.RegisteredClaims
}

func (j *JwtWrapper) ParseToken(tokenString string) (claims *JwtCustomClaim, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtCustomClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JwtCustomClaim)
	if !ok {
		return nil, fmt.Errorf("claims not JwtClaim")
	}

	return claims, nil
}

func (j *JwtWrapper) ValidateToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtCustomClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return false
	}

	claims, _ := token.Claims.(*JwtCustomClaim)
	// if !cache.UserSessionCheck(claims.UserID, tokenString) {
	// 	return false
	// }

	if claims.ExpiresAt.Local().Unix() < time.Now().Local().Unix() {
		return false
	}

	return token.Valid
}

func (j *JwtWrapper) GenerateJWT(userName, userId string) (string, error) {
	claims := &JwtClaim{
		UserName: userName,
		UserId:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(config.ReadValue().App.JwtExpire))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
