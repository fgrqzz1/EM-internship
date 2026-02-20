package models

import (
	"time"
)

type Subscription struct {
	ID          string    `json:"id" db:"id"`
	ServiceName string    `json:"service_name" db:"service_name" validate:"required,min=1,max=255"`
	Price       int       `json:"price" db:"price" validate:"required"`
	UserID      string    `json:"user_id" db:"user_id" validate:"required,uuid"`
	StartDate   string    `json:"start_date" db:"start_date" validate:"required"`
	EndDate     *string   `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSubscriptionInput struct {
	ServiceName string `json:"service_name" validate:"required,min=1,max=255"`
	Price       int    `json:"price" validate:"required"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date,omitempty"`
}

type UpdateSubscriptionInput struct {
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=255"`
	Price       *int    `json:"price,omitempty" validate:"omitempty"`
	UserID      *string `json:"user_id,omitempty" validate:"omitempty,uuid"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
	Count     int `json:"count"`
}

type SubscriptionList struct {
	Items []Subscription `json:"items"`
	Total int            `json:"total"`
}
