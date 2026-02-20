package repository

import (
	"context"
	"database/sql"
	"em-internship/internal/models"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewSubscriptionRepository(db *pgxpool.Pool, logger *zap.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, input models.CreateSubscriptionInput) (*models.Subscription, error) {
	id := uuid.New().String()
	nowTime := time.Now()

	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var sub models.Subscription
	var endDate sql.NullString
	if input.EndDate != "" {
		endDate = sql.NullString{String: input.EndDate, Valid: true}
	}

	err := r.db.QueryRow(ctx, query,
		id, input.ServiceName, input.Price, input.UserID, input.StartDate, endDate, nowTime, nowTime,
	).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt)

	if err != nil {
		r.logger.Error("failed to create subscription", zap.Error(err), zap.String("user_id", input.UserID))
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	r.logger.Info("created subscription", zap.String("id", sub.ID), zap.String("user_id", input.UserID))

	return &sub, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var sub models.Subscription
	var endDate sql.NullString

	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrSubscriptionNotFound
	}

	if err != nil {
		r.logger.Error("failed to get subscription", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	if endDate.Valid {
		sub.EndDate = &endDate.String
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetAll(ctx context.Context, limit, offset int) (*models.SubscriptionList, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	countQuery := "SELECT COUNT(*) FROM subscriptions"

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to get subscriptions", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		var endDate sql.NullString

		err := rows.Scan(
			&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
			&sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			sub.EndDate = &endDate.String
		}

		subs = append(subs, sub)
	}

	var total int
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		r.logger.Error("failed to count subscriptions", zap.Error(err))
	}

	return &models.SubscriptionList{
		Items: subs,
		Total: total,
	}, nil
}
