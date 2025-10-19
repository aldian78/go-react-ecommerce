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

type ICategoryRepository interface {
	CreateNewCategory(ctx context.Context, Category *entity.Category) error
	GetCategoryById(ctx context.Context, id string) (*entity.Category, error)
	GetCategoryByIds(ctx context.Context, ids []string) ([]*entity.Category, error)
	UpdateCategory(ctx context.Context, Category *entity.Category) error
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryPagination(ctx context.Context, pagination *basecommon.PaginationRequest) ([]*entity.Category, *basecommon.PaginationResponse, error)
}

type categoryRepository struct {
	db database.DatabaseQuery
}

func (repo *categoryRepository) CreateNewCategory(ctx context.Context, Category *entity.Category) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO category (id, name, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6)",
		Category.Id,
		Category.Name,
		Category.CreatedAt,
		Category.CreatedBy,
		Category.UpdatedAt,
		Category.UpdatedBy,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *categoryRepository) GetCategoryById(ctx context.Context, id string) (*entity.Category, error) {
	var CategoryEntity entity.Category
	row := repo.db.QueryRowContext(
		ctx,
		"SELECT id, name FROM category WHERE id = $1",
		id,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&CategoryEntity.Id,
		&CategoryEntity.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &CategoryEntity, nil
}

func (repo *categoryRepository) GetCategoryByIds(ctx context.Context, ids []string) ([]*entity.Category, error) {
	queryIds := make([]string, len(ids))
	for i, id := range ids {
		queryIds[i] = fmt.Sprintf("'%s'", id)
	}
	rows, err := repo.db.QueryContext(
		ctx,
		fmt.Sprintf("SELECT id, name FROM category WHERE id IN (%s)", strings.Join(queryIds, ", ")),
	)
	if err != nil {
		return nil, err
	}

	var Category []*entity.Category = make([]*entity.Category, 0)
	for rows.Next() {
		var CategoryEntity entity.Category
		err = rows.Scan(
			&CategoryEntity.Id,
			&CategoryEntity.Name,
		)
		if err != nil {
			return nil, err
		}

		Category = append(Category, &CategoryEntity)
	}

	return Category, nil
}

func (repo *categoryRepository) UpdateCategory(ctx context.Context, Category *entity.Category) error {
	_, err := repo.db.ExecContext(
		ctx,
		"UPDATE category SET name=$1, updated_at=$2, updated_by=$3 WHERE id = $4",
		Category.Name,
		Category.UpdatedAt,
		Category.UpdatedBy,
		Category.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *categoryRepository) DeleteCategory(ctx context.Context, id string) error {
	_, err := repo.db.ExecContext(
		ctx,
		"DELETE FROM category WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *categoryRepository) GetCategoryPagination(ctx context.Context, pagination *basecommon.PaginationRequest) ([]*entity.Category, *basecommon.PaginationResponse, error) {

	row := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Category")
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
		"SELECT id, name FROM Category ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		pagination.ItemPerPage,
		offset,
	)
	if err != nil {
		return nil, nil, err
	}

	var category []*entity.Category = make([]*entity.Category, 0)
	for rows.Next() {
		var categorys entity.Category

		err = rows.Scan(
			&categorys.Id,
			&categorys.Name,
		)
		if err != nil {
			return nil, nil, err
		}

		category = append(category, &categorys)
	}

	paginationResponse := &basecommon.PaginationResponse{
		CurrentPage:    pagination.CurrentPage,
		ItemPerPage:    pagination.ItemPerPage,
		TotalItemCount: int32(totalCount),
		TotalPageCount: int32(totalPages),
	}
	return category, paginationResponse, nil
}

func NewCategoryRepository(db database.DatabaseQuery) ICategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
