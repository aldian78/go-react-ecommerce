package jwt

import (
	"context"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/common/utils"
	gocache "github.com/patrickmn/go-cache"
	"go-micro.dev/v4/logger"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtEntityContextKey string

var JwtEntityContextKeyValue JwtEntityContextKey = "JwtEntity"

type JwtClaims struct {
	jwt.RegisteredClaims
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func (jc *JwtClaims) SetToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, JwtEntityContextKeyValue, jc)
}

func GetClaimsFromToken(token string) (*JwtClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	if !tokenClaims.Valid {
		return nil, utils.UnauthenticatedResponse()
	}

	if claims, ok := tokenClaims.Claims.(*JwtClaims); ok {
		return claims, nil
	}

	return nil, utils.UnauthenticatedResponse()
}

func GetClaimsFromContext(token string) (*JwtClaims, error) {
	logger.Infof("check ctrsss : %s", token)
	//claims, ok := ctx.Value(JwtEntityContextKeyValue).(*JwtClaims)
	//logger.Infof("check claims : %", claims)
	//if !ok {
	//	return nil, utils.UnauthenticatedResponse()
	//}

	cacheService := gocache.New(time.Hour*24, time.Hour)
	tokens, ok := cacheService.Get(token)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}
	claims := tokens.(*JwtClaims)
	logger.Infof("check final claims : %", claims)
	return claims, nil
}
