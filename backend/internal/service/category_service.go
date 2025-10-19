package service

import (
	"context"
	"github.com/aldian78/go-react-ecommerce/backend/internal/entity"
	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	baseutil "github.com/aldian78/go-react-ecommerce/common/utils"
	"github.com/aldian78/go-react-ecommerce/proto/pb/category"
	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
	"time"
)

type ICategoryService interface {
	CreateCategory(ctx context.Context, request *category.CreateCategoryRequest, params map[string]string) (*category.CreateCategoryResponse, error)
	EditCategory(ctx context.Context, request *category.EditCategoryRequest, params map[string]string) (*category.EditCategoryResponse, error)
	DeleteCategory(ctx context.Context, request *category.DeleteCategoryRequest, params map[string]string) (*category.DeleteCategoryResponse, error)
	ListCategory(ctx context.Context, request *category.ListCategoryRequest) (*category.ListCategoryResponse, error)
}

type categoryService struct {
	categoryRepository repository.ICategoryRepository
}

func (ps *categoryService) CreateCategory(ctx context.Context, request *category.CreateCategoryRequest, params map[string]string) (*category.CreateCategoryResponse, error) {
	if params["role"] != entity.UserRoleAdmin {
		return nil, baseutil.UnauthenticatedResponse()
	}

	categoryEntity := entity.Category{
		Id:        uuid.NewString(),
		Name:      request.Name,
		CreatedAt: time.Now(),
		CreatedBy: params["fullName"],
	}
	err := ps.categoryRepository.CreateNewCategory(ctx, &categoryEntity)
	if err != nil {
		return nil, err
	}

	return &category.CreateCategoryResponse{
		Base: baseutil.SuccessResponse("Category is created"),
		Id:   categoryEntity.Id,
	}, nil
}

func (ps *categoryService) EditCategory(ctx context.Context, request *category.EditCategoryRequest, params map[string]string) (*category.EditCategoryResponse, error) {
	if params["role"] != entity.UserRoleAdmin {
		return nil, baseutil.UnauthenticatedResponse()
	}

	categoryEntity, err := ps.categoryRepository.GetCategoryById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if categoryEntity == nil {
		return &category.EditCategoryResponse{
			Base: baseutil.NotFoundResponse("Category not found"),
		}, nil
	}

	fullName := params["fullName"]
	newProduct := entity.Category{
		Id:        request.Id,
		Name:      request.Name,
		UpdatedAt: time.Now(),
		UpdatedBy: &fullName,
	}

	logger.Infof("check new product : %s", &newProduct)

	err = ps.categoryRepository.UpdateCategory(ctx, &newProduct)
	if err != nil {
		return nil, err
	}

	return &category.EditCategoryResponse{
		Base: baseutil.SuccessResponse("Updated Category success"),
		Id:   request.Id,
	}, nil
}

func (ps *categoryService) DeleteCategory(ctx context.Context, request *category.DeleteCategoryRequest, params map[string]string) (*category.DeleteCategoryResponse, error) {
	logger.Infof("check role : %s", params["role"])
	if params["role"] != entity.UserRoleAdmin {
		return nil, baseutil.UnauthenticatedResponse()
	}

	productEntity, err := ps.categoryRepository.GetCategoryById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if productEntity == nil {
		return &category.DeleteCategoryResponse{
			Base: baseutil.NotFoundResponse("Category not found"),
		}, nil
	}

	err = ps.categoryRepository.DeleteCategory(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return &category.DeleteCategoryResponse{
		Base: baseutil.SuccessResponse("Deleted Category success"),
	}, nil
}

func (ps *categoryService) ListCategory(ctx context.Context, request *category.ListCategoryRequest) (*category.ListCategoryResponse, error) {
	categories, paginationResponse, err := ps.categoryRepository.GetCategoryPagination(ctx, request.Pagination)
	if err != nil {
		return nil, err
	}

	var data []*category.ListCategoryResponseItem = make([]*category.ListCategoryResponseItem, 0)
	for _, prod := range categories {
		data = append(data, &category.ListCategoryResponseItem{
			Id:   prod.Id,
			Name: prod.Name,
		})
	}

	return &category.ListCategoryResponse{
		Base:       baseutil.SuccessResponse("Get list category success"),
		Pagination: paginationResponse,
		Data:       data,
	}, nil
}

func NewCategoryService(categoryRepository repository.ICategoryRepository) ICategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
