package entity

import "time"

type ProductType struct {
	Id        string
	Type      string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy *string
}
