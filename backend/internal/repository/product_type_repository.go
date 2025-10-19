package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/backend/internal/entity"
	"github.com/aldian78/go-react-ecommerce/proto/pb/basecommon"
	"strings"

	"github.com/aldian78/go-react-ecommerce/backend/pkg/database"
)

type IProductTypeRepository interface {
	CreateNewProductType(ctx context.Context, ProductType *entity.ProductType) error
	GetProductTypeById(ctx context.Context, id string) (*entity.ProductType, error)
	GetProductTypeByIds(ctx context.Context, ids []string) ([]*entity.ProductType, error)
	UpdateProductType(ctx context.Context, ProductType *entity.ProductType) error
	DeleteProductType(ctx context.Context, id string) error
	GetProductTypePagination(ctx context.Context, pagination *basecommon.PaginationRequest) ([]*entity.ProductType, *basecommon.PaginationResponse, error)
}

type productTypeRepository struct {
	db database.DatabaseQuery
}

func (repo *productTypeRepository) CreateNewProductType(ctx context.Context, ProductType *entity.ProductType) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO product_type (id, type, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6)",
		ProductType.Id,
		ProductType.Type,
		ProductType.CreatedAt,
		ProductType.CreatedBy,
		ProductType.UpdatedAt,
		ProductType.UpdatedBy,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *productTypeRepository) GetProductTypeById(ctx context.Context, id string) (*entity.ProductType, error) {
	var ProductTypeEntity entity.ProductType
	row := repo.db.QueryRowContext(
		ctx,
		"SELECT id, type FROM product_type WHERE id = $1",
		id,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&ProductTypeEntity.Id,
		&ProductTypeEntity.Type,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &ProductTypeEntity, nil
}

func (repo *productTypeRepository) GetProductTypeByIds(ctx context.Context, ids []string) ([]*entity.ProductType, error) {
	queryIds := make([]string, len(ids))
	for i, id := range ids {
		queryIds[i] = fmt.Sprintf("'%s'", id)
	}
	rows, err := repo.db.QueryContext(
		ctx,
		fmt.Sprintf("SELECT id, type FROM product_type WHERE id IN (%s)", strings.Join(queryIds, ", ")),
	)
	if err != nil {
		return nil, err
	}

	var ProductType []*entity.ProductType = make([]*entity.ProductType, 0)
	for rows.Next() {
		var ProductTypeEntity entity.ProductType
		err = rows.Scan(
			&ProductTypeEntity.Id,
			&ProductTypeEntity.Type,
		)
		if err != nil {
			return nil, err
		}

		ProductType = append(ProductType, &ProductTypeEntity)
	}

	return ProductType, nil
}

func (repo *productTypeRepository) UpdateProductType(ctx context.Context, ProductType *entity.ProductType) error {
	_, err := repo.db.ExecContext(
		ctx,
		"UPDATE product_type SET type=$1, updated_at=$2, updated_by=$3 WHERE id = $4",
		ProductType.Type,
		ProductType.UpdatedAt,
		ProductType.UpdatedBy,
		ProductType.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *productTypeRepository) DeleteProductType(ctx context.Context, id string) error {
	_, err := repo.db.ExecContext(
		ctx,
		"DELETE FROM product_type WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *productTypeRepository) GetProductTypePagination(ctx context.Context, pagination *basecommon.PaginationRequest) ([]*entity.ProductType, *basecommon.PaginationResponse, error) {

	row := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM product_type")
	if row.Err() != nil {
		return nil, nil, row.Err()
	}

	var totalCount int
	err := row.Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	offset := (pagination.CurrentPage - 1) * pagination.ItemPerPage
	totalPages := (totalCount + int(pagination.ItemPerPage) - 1) / int(pagination.ItemPerPage)

	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, type FROM product_type ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		pagination.ItemPerPage,
		offset,
	)
	if err != nil {
		return nil, nil, err
	}

	var ProductType []*entity.ProductType = make([]*entity.ProductType, 0)
	for rows.Next() {
		var ProductTypes entity.ProductType

		err = rows.Scan(
			&ProductTypes.Id,
			&ProductTypes.Type,
		)
		if err != nil {
			return nil, nil, err
		}

		ProductType = append(ProductType, &ProductTypes)
	}

	paginationResponse := &basecommon.PaginationResponse{
		CurrentPage:    pagination.CurrentPage,
		ItemPerPage:    pagination.ItemPerPage,
		TotalItemCount: int32(totalCount),
		TotalPageCount: int32(totalPages),
	}
	return ProductType, paginationResponse, nil
}

func NewProductTypeRepository(db database.DatabaseQuery) IProductTypeRepository {
	return &productTypeRepository{
		db: db,
	}
}
