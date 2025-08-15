package entity

import (
	"time"
)

const (
	Routes = "routes"
)

type ApiGateway struct {
	Id           string
	TxType       string
	ContextPath  string
	Method       string
	Middleware   []string
	EndpointName string
	EndpointFunc string
	ContentType  string
	Description  string
	CreatedBy    *string
	CreatedAt    time.Time
	UpdatedBy    *string
	UpdatedAt    time.Time
}
