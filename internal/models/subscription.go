package models

import (
	"time"
)

// Subscription модель подписки
// @Description Модель подписки на сервис
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

// для создания подписки
type CreateSubscriptionInput struct {
	ServiceName string `json:"service_name" validate:"required,min=1,max=255"`
	Price       int    `json:"price" validate:"required,gt=0"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	StartDate   string `json:"start_date" validate:"required,month_year"`
	EndDate     string `json:"end_date,omitempty" validate:"omitempty,month_year"`
}

// для обновления подписки
type UpdateSubscriptionInput struct {
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=255"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,gt=0"`
	UserID      *string `json:"user_id,omitempty" validate:"omitempty,uuid"`
	StartDate   *string `json:"start_date,omitempty" validate:"omitempty,month_year"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,month_year"`
}

// итоговая стоимость
type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
	Count     int `json:"count"`
}

// список всех подписок
type SubscriptionList struct {
	Items []Subscription `json:"items"`
	Total int            `json:"total"`
}
