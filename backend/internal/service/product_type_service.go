package service

import (
	"context"
	"github.com/aldian78/go-react-ecommerce/backend/internal/entity"
	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	baseutil "github.com/aldian78/go-react-ecommerce/common/utils"
	ProductType "github.com/aldian78/go-react-ecommerce/proto/pb/product_type"
	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
	"time"
)

type IProductTypeService interface {
	CreateProductType(ctx context.Context, request *ProductType.CreateProductTypeRequest, params map[string]string) (*ProductType.CreateProductTypeResponse, error)
	EditProductType(ctx context.Context, request *ProductType.EditProductTypeRequest, params map[string]string) (*ProductType.EditProductTypeResponse, error)
	DeleteProductType(ctx context.Context, request *ProductType.DeleteProductTypeRequest, params map[string]string) (*ProductType.DeleteProductTypeResponse, error)
	ListProductType(ctx context.Context, request *ProductType.ListProductTypeRequest) (*ProductType.ListProductTypeResponse, error)
}

type productTypeService struct {
	productTypeRepository repository.IProductTypeRepository
}

func (ps *productTypeService) CreateProductType(ctx context.Context, request *ProductType.CreateProductTypeRequest, params map[string]string) (*ProductType.CreateProductTypeResponse, error) {
	if params["role"] != entity.UserRoleAdmin {
		return nil, baseutil.UnauthenticatedResponse()
	}

	ProductTypeEntity := entity.ProductType{
		Id:        uuid.NewString(),
		Type:      request.Type,
		CreatedAt: time.Now(),
		CreatedBy: params["fullName"],
	}
	err := ps.productTypeRepository.CreateNewProductType(ctx, &ProductTypeEntity)
	if err != nil {
		return nil, err
	}

	return &ProductType.CreateProductTypeResponse{
		Base: baseutil.SuccessResponse("ProductType is created"),
		Id:   ProductTypeEntity.Id,
	}, nil
}

func (ps *productTypeService) EditProductType(ctx context.Context, request *ProductType.EditProductTypeRequest, params map[string]string) (*ProductType.EditProductTypeResponse, error) {
	if params["role"] != entity.UserRoleAdmin {
		return nil, baseutil.UnauthenticatedResponse()
	}

	ProductTypeEntity, err := ps.productTypeRepository.GetProductTypeById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if ProductTypeEntity == nil {
		return &ProductType.EditProductTypeResponse{
			Base: baseutil.NotFoundResponse("ProductType not found"),
		}, nil
	}

	fullName := params["fullName"]
	newProduct := entity.ProductType{
		Id:        request.Id,
		Type:      request.Type,
		UpdatedAt: time.Now(),
		UpdatedBy: &fullName,
	}

	logger.Infof("check new product : %s", &newProduct)

	err = ps.productTypeRepository.UpdateProductType(ctx, &newProduct)
	if err != nil {
		return nil, err
	}

	return &ProductType.EditProductTypeResponse{
		Base: baseutil.SuccessResponse("Updated ProductType success"),
		Id:   request.Id,
	}, nil
}

func (ps *productTypeService) DeleteProductType(ctx context.Context, request *ProductType.DeleteProductTypeRequest, params map[string]string) (*ProductType.DeleteProductTypeResponse, error) {
	logger.Infof("check role : %s", params["role"])
	if params["role"] != entity.UserRoleAdmin {
		return nil, baseutil.UnauthenticatedResponse()
	}

	productEntity, err := ps.productTypeRepository.GetProductTypeById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if productEntity == nil {
		return &ProductType.DeleteProductTypeResponse{
			Base: baseutil.NotFoundResponse("ProductType not found"),
		}, nil
	}

	err = ps.productTypeRepository.DeleteProductType(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &ProductType.DeleteProductTypeResponse{
		Base: baseutil.SuccessResponse("Deleted ProductType success"),
	}, nil
}

func (ps *productTypeService) ListProductType(ctx context.Context, request *ProductType.ListProductTypeRequest) (*ProductType.ListProductTypeResponse, error) {
	productType, paginationResponse, err := ps.productTypeRepository.GetProductTypePagination(ctx, request.Pagination)
	if err != nil {
		return nil, err
	}

	var data []*ProductType.ListProductTypeResponseItem = make([]*ProductType.ListProductTypeResponseItem, 0)
	for _, prod := range productType {
		data = append(data, &ProductType.ListProductTypeResponseItem{
			Id:   prod.Id,
			Type: prod.Type,
		})
	}

	return &ProductType.ListProductTypeResponse{
		Base:       baseutil.SuccessResponse("Get list ProductType success"),
		Pagination: paginationResponse,
		Data:       data,
	}, nil
}

func NewProductTypeService(productTypeRepository repository.IProductTypeRepository) IProductTypeService {
	return &productTypeService{
		productTypeRepository: productTypeRepository,
	}
}
