package main

import (
	"context"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/handler"
	"github.com/aldian78/go-react-ecommerce/backend/pkg/database"
	auth "github.com/aldian78/go-react-ecommerce/proto/pb/authentication"
	"github.com/aldian78/go-react-ecommerce/proto/pb/cart"
	"github.com/aldian78/go-react-ecommerce/proto/pb/newsletter"
	"github.com/aldian78/go-react-ecommerce/proto/pb/order"
	"github.com/aldian78/go-react-ecommerce/proto/pb/product"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	gtc "github.com/shengyanli1982/go-trycatch"
	"github.com/xendit/xendit-go"
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

			xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET")
			logger.Infof("check xendit secret : %s", os.Getenv("XENDIT_SECRET"))

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

			addrRedis := os.Getenv("USER_REDIS")

			logger.Infof("redis addr: %v", addrRedis)
			//passRedis := os.Getenv("PASS_REDIS")
			rdb := redis.NewClient(&redis.Options{
				Addr:     addrRedis,
				Password: "",
				DB:       0, // default DB
			})

			_ = svc.Client().Init(client.RequestTimeout(rpcTimeout))
			_ = auth.RegisterAuthenticationServiceHandler(svc.Server(), handler.NewAuthenticationHandler(db, rdb, cacheService))
			_ = product.RegisterProductServiceHandler(svc.Server(), handler.NewProductHandler(db))
			_ = cart.RegisterCartServiceHandler(svc.Server(), handler.NewCartHandler(db))
			_ = newsletter.RegisterNewsletterServiceHandler(svc.Server(), handler.NewNewsletterHandler(db))
			_ = order.RegisterOrderServiceHandler(svc.Server(), handler.NewOrderHandler(db))

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
