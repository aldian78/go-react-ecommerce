package service

import (
	"context"
	"github.com/aldian78/go-react-ecommerce/backend/internal/entity"
	baseutil "github.com/aldian78/go-react-ecommerce/common/utils"
	"time"

	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	"github.com/aldian78/go-react-ecommerce/proto/pb/newsletter"
	"github.com/google/uuid"
)

type INewsletterService interface {
	SubscribeNewsletter(ctx context.Context, request *newsletter.SubcribeNewsletterRequest) (*newsletter.SubcribeNewsletterResponse, error)
}

type newsletterService struct {
	newsletterRepository repository.INewsletterRepository
}

func (ns *newsletterService) SubscribeNewsletter(ctx context.Context, request *newsletter.SubcribeNewsletterRequest) (*newsletter.SubcribeNewsletterResponse, error) {
	newsletterEntity, err := ns.newsletterRepository.GetNewsletterByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if newsletterEntity != nil {
		return &newsletter.SubcribeNewsletterResponse{
			Base: baseutil.SuccessResponse("Subscribe newsletter success"),
		}, nil
	}

	newNewsletterEntity := entity.Newsletter{
		Id:        uuid.NewString(),
		FullName:  request.FullName,
		Email:     request.Email,
		CreatedAt: time.Now(),
		CreatedBy: "Public",
	}
	err = ns.newsletterRepository.CreateNewNewsletter(ctx, &newNewsletterEntity)
	if err != nil {
		return nil, err
	}

	return &newsletter.SubcribeNewsletterResponse{
		Base: baseutil.SuccessResponse("Subscribe newsletter success"),
	}, nil
}

func NewNewsletterService(newsletterRepository repository.INewsletterRepository) INewsletterService {
	return &newsletterService{
		newsletterRepository: newsletterRepository,
	}
}
