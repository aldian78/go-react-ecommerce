package grpcmiddleware

import (
	"context"
	"github.com/aldian78/go-react-ecommerce/common/jwt"

	"github.com/aldian78/go-react-ecommerce/backend/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

var publicApis = map[string]bool{
	"/auth.AuthService/Login":                           true,
	"/auth.AuthService/Register":                        true,
	"/product.ProductService/DetailProduct":             true,
	"/product.ProductService/ListProduct":               true,
	"/product.ProductService/HighlightProducts":         true,
	"/newsletter.NewsletterService/SubscribeNewsletter": true,
}

func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if publicApis[info.FullMethod] {
		return handler(ctx, req)
	}

	tokenStr, err := jwt.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, ok := am.cacheService.Get(tokenStr)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	claims, err := jwt.GetClaimsFromToken(tokenStr)
	if err != nil {
		return nil, err
	}

	ctx = claims.SetToContext(ctx)

	res, err := handler(ctx, req)

	return res, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}
