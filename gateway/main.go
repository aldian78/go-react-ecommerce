package main

import (
	"context"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/gateway/internal/repository"
	"github.com/aldian78/go-react-ecommerce/gateway/pkg/database"
	"github.com/aldian78/go-react-ecommerce/gateway/service"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/web"
	"github.com/go-micro/plugins/v3/client/grpc"
	"github.com/joho/godotenv"
	gtc "github.com/shengyanli1982/go-trycatch"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	gtc.New().
		Try(func() error {
			svc := web.NewService(
				web.Name("proto-api"),
				web.Address(":8081"),
			)

			godotenv.Load()

			ctx := context.Background()
			db := database.ConnectDB(ctx, os.Getenv("DB_URI"))
			routerRepository := repository.NewRoutesRepository(db)

			gwTimeout, _ := strconv.Atoi("60")
			if gwTimeout <= 0 {
				gwTimeout = 60
			}

			// Membuat custom gRPC client
			_ = grpc.NewClient(
				client.RequestTimeout(time.Duration(gwTimeout)*time.Second),
				grpc.MaxRecvMsgSize(15*1024*1024),
				grpc.MaxSendMsgSize(15*1024*1024),
			)

			corsOrigins := "http://127.0.0.1:8080"
			corsMethods := "GET,POST,PUT,DELETE,OPTIONS"
			corsHeaders := "Content-Type,Authorization,X-Requested-With"

			corsOrigin := strings.Split(corsOrigins, ",")
			corsMethod := strings.Split(corsMethods, ",")
			corsHeader := strings.Split(corsHeaders, ",")

			//services.InitService()

			svc.Handle("/", service.Init("api", gwTimeout, corsOrigin, corsMethod, corsHeader, routerRepository))
			if err := svc.Run(); err != nil {
				panic(err)
			}
			return fmt.Errorf("this try in main.go")
		}).
		Catch(func(err error) {
			fmt.Printf("Caught error: %v\n", err)
		}).
		Do()
}
