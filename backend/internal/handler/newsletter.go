package handler

import (
	"context"
	baseutil "github.com/aldian78/go-react-ecommerce/common/utils"

	"github.com/aldian78/go-react-ecommerce/backend/internal/service"
	"github.com/aldian78/go-react-ecommerce/backend/internal/utils"
	"github.com/aldian78/go-react-ecommerce/proto/pb/newsletter"
)

type newsletterHandler struct {
	newsletter.UnimplementedNewsletterServiceServer

	newsletterService service.INewsletterService
}

func (nh *newsletterHandler) SubscribeNewsletter(ctx context.Context, request *newsletter.SubcribeNewsletterRequest) (*newsletter.SubcribeNewsletterResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &newsletter.SubcribeNewsletterResponse{
			Base: baseutil.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := nh.newsletterService.SubscribeNewsletter(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewNewsletterHandler(newsletterService service.INewsletterService) *newsletterHandler {
	return &newsletterHandler{
		newsletterService: newsletterService,
	}
}
