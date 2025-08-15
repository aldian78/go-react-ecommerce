package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
	gtc "github.com/shengyanli1982/go-trycatch"
	"go-grpc-ecommerce-be/internal/handler"
	auth "go-grpc-ecommerce-be/pb/authentication"
	"go-grpc-ecommerce-be/pkg/database"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {

	gtc.New().
		Try(func() error {
			runtime.GOMAXPROCS(4)

			// Define services
			svc := micro.NewService(
				micro.Name("backend"),
				micro.Address(":8085"),
			)

			godotenv.Load()

			rto, _ := strconv.Atoi("60")
			var rpcTimeout = time.Second * 60
			if rto > 0 {
				rpcTimeout = time.Duration(rto) * time.Second
			}
			_ = svc.Client().Init(client.RequestTimeout(rpcTimeout))

			logger.Info("services timeout in " + rpcTimeout.String())

			ctx := context.Background()
			db := database.ConnectDB(ctx, os.Getenv("DB_URI"))

			cacheService := gocache.New(time.Hour*24, time.Hour)

			_ = svc.Client().Init(client.RequestTimeout(rpcTimeout))
			_ = auth.RegisterAuthenticationServiceHandler(svc.Server(), handler.NewAuthenticationHandler(db, cacheService))

			// Run Service
			if err := svc.Run(); err != nil {
				logger.Error(err)
				panic(err)
			}
			return fmt.Errorf("this try in services.go")
		}).
		Catch(func(err error) {
			fmt.Printf("Caught error: %v\n", err)
		}).
		Do()
}
