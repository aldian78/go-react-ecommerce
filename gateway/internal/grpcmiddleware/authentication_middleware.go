package grpcmiddleware

import (
	"context"
	jwtentity "github.com/aldian78/go-react-ecommerce/backend/pkg/entity/jwt"
	"github.com/gin-gonic/gin"
	gocache "github.com/patrickmn/go-cache"
	"go-micro.dev/v4/logger"
	"google.golang.org/grpc/metadata"
	"net/http"
	"time"
)

//
//func AuthenticationMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token := c.GetHeader("Authorization")
//		if token == "" {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
//			return
//		}
//
//		logger.Infof("check token: %v", token)
//
//		// Validasi token sama seperti di gRPC middleware
//		cacheService := gocache.New(time.Hour*24, time.Hour)
//		if _, found := cacheService.Get(token); found {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
//			return
//		}
//
//		claims, err := jwtentity.GetClaimsFromToken(strings.TrimPrefix(token, "Bearer "))
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
//			return
//		}
//
//		logger.Infof("check res claims: %v", claims)
//
//		// Simpan claims ke Gin context
//		c.Set("JwtEntity", claims)
//
//		c.Next()
//	}
//}

func NewAuthenticationMiddleware() {
	AuthorizationMiddlewareFactory["AuthenticationMiddleware"] = AuthenticationMiddleware()
}

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.TODO()
		authHeader := c.GetHeader("Authorization")
		logger.Infof("Auth Header: %s", authHeader)

		md := metadata.Pairs("authorization", authHeader)
		ctx = metadata.NewIncomingContext(ctx, md)

		logger.Info("ctx :", ctx)
		tokenStr, err := jwtentity.ParseTokenFromContext(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked 1"})
			return
		}

		cacheService := gocache.New(time.Hour*24, time.Hour)
		_, ok := cacheService.Get(tokenStr)
		if ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked 2"})
			return
		}

		claims, err := jwtentity.GetClaimsFromToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked 3"})
			return
		}

		ctx = claims.SetToContext(ctx)

		// Simpan context baru ini ke Gin supaya bisa dipakai di handler
		c.Request = c.Request.WithContext(ctx)

		logger.Info("after claims set to context :", c)

		c.Next()
	}
}
