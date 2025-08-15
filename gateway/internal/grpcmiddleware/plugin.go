package grpcmiddleware

import (
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"path/filepath"
	"strings"
)

var (
	baseDir = "/app/middleware"

	middlewares = make(map[string]gin.HandlerFunc)
)

func LoadMiddlewareFromFactory(name []string) []gin.HandlerFunc {
	handler := make([]gin.HandlerFunc, 0)

	for _, n := range name {
		filename := strings.TrimSuffix(n, filepath.Ext(n))
		logger.Infof("Loading middleware from file %s", filename)
		h := AuthorizationMiddlewareFactory[filename]
		handler = append(handler, h)
	}

	return handler
}
