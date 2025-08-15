package grpcmiddleware

import "github.com/gin-gonic/gin"

var AuthorizationMiddlewareFactory = make(map[string]gin.HandlerFunc)

func InitAuthorizationMiddleware() {
	NewAuthenticationMiddleware()
}
