package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go-grpc-ecommerce-be/internal/entity"
	"go-grpc-ecommerce-be/internal/repository"
	"go-grpc-ecommerce-be/internal/utils"
	"go-grpc-ecommerce-be/pb/newsletter"
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
			Base: utils.SuccessResponse("Subscribe newsletter success"),
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
		Base: utils.SuccessResponse("Subscribe newsletter success"),
	}, nil
}

func NewNewsletterService(newsletterRepository repository.INewsletterRepository) INewsletterService {
	return &newsletterService{
		newsletterRepository: newsletterRepository,
	}
}
