package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-micro.dev/v4/logger"
	"strings"

	"github.com/aldian78/go-react-ecommerce/common/utils"
	"google.golang.org/grpc/metadata"
)

func ParseTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	logger.Infof("check md 1 : %", md)
	if !ok {
		logger.Infof("check ok : %", ok)
		return "", utils.UnauthenticatedResponse()
	}

	bearerToken, ok := md["authorization"]
	logger.Infof("check bearerToken : %", bearerToken)
	logger.Infof("check ok 2 : %", ok)
	if !ok {
		return "", utils.UnauthenticatedResponse()
	}

	if len(bearerToken) == 0 {
		return "", utils.UnauthenticatedResponse()
	}

	tokenSplit := strings.Split(bearerToken[0], " ")

	logger.Infof("check tokenSplit : %", tokenSplit)

	if len(tokenSplit) != 2 {
		return "", utils.UnauthenticatedResponse()
	}

	logger.Infof("check tokenSplit[0] : %", tokenSplit[0])

	if tokenSplit[0] != "Bearer" {
		return "", utils.UnauthenticatedResponse()
	}

	logger.Infof("check tokenSplit[1] : %", tokenSplit[1])
	return tokenSplit[1], nil
}

func ParseToken(tokenStr string) (string, error) {
	if tokenStr == "" {
		return "", fmt.Errorf("missing token")
	}

	// Kalau ada prefix "Bearer " dihilangkan
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = strings.TrimSpace(tokenStr[7:])
	}

	logger.Info("check tokenStr 2 : ", tokenStr)
	// Validasi format dasar JWT (3 bagian dipisahkan ".")
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	return tokenStr, nil
}

func ParseTokenJWT(tokenStr string, secretKey string) (*JwtClaims, error) {
	claims := &JwtClaims{}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Pastikan algoritma HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Validasi
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
