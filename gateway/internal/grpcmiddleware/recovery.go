package grpcmiddleware

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func GetConfig() map[string]string {
	return make(map[string]string)
}

// Recovery is error handling middleware
func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool

			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			logger.Infof("internal server error: %v", err, string(debug.Stack()))

			// If the connection is dead, we can't write a status to it.
			if brokenPipe {
				c.Error(err.(error)) // nolint: errcheck
				c.Abort()
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()

	// Add header cache-control
	if len(GetConfig()["CACHE_CONTROL_HEADER"]) != 0 {
		c.Header("Cache-Control", GetConfig()["CACHE_CONTROL_HEADER"])
	}

	c.Next()
}
