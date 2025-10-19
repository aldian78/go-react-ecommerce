package entity

import "time"

type Category struct {
	Id        string
	Name      string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy *string
}
