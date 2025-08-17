package repository

import (
	"context"
	"fmt"
	"github.com/aldian78/go-react-ecommerce/gateway/internal/entity"
	"github.com/aldian78/go-react-ecommerce/gateway/pkg/database"
	"github.com/lib/pq"
	"log"
)

type IRoutesRepository interface {
	GetAllEndpoint(ctx context.Context) ([]*entity.ApiGateway, error)
}
type routesRepository struct {
	db database.DatabaseQuery
}

func (ar *routesRepository) GetAllEndpoint(ctx context.Context) ([]*entity.ApiGateway, error) {
	row, err := ar.db.QueryContext(ctx, "SELECT * FROM routes")
	if err != nil {
		fmt.Println("errror su ", err.Error())
		return nil, row.Err()
	}
	log.Printf("row : %s", row)
	var apis []*entity.ApiGateway = make([]*entity.ApiGateway, 0)

	defer row.Close()
	for row.Next() {
		var api entity.ApiGateway
		err := row.Scan(
			&api.Id,
			&api.TxType,
			&api.ContextPath,
			&api.Method,
			pq.Array(&api.Middleware),
			&api.EndpointName,
			&api.EndpointFunc,
			&api.ContentType,
			&api.Description,
			&api.CreatedBy,
			&api.CreatedAt,
			&api.UpdatedBy,
			&api.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		log.Printf("proto-api : [%s] %s", &api)
		apis = append(apis, &api)
	}

	return apis, nil
}

func NewRoutesRepository(db database.DatabaseQuery) IRoutesRepository {
	return &routesRepository{
		db: db,
	}
}
