package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"em-internship/internal/models"
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

func (r *SubscriptionRepository) Update(ctx context.Context, id string, input models.UpdateSubscriptionInput) (*models.Subscription, error) {
	query := `
		UPDATE subscriptions 
		SET service_name = COALESCE($1, service_name),
			price = COALESCE($2, price),
			user_id = COALESCE($3, user_id),
			start_date = COALESCE($4, start_date),
			end_date = COALESCE($5, end_date),
			updated_at = $6
		WHERE id = $7
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	sub, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	serviceName := sub.ServiceName
	if input.ServiceName != nil {
		serviceName = *input.ServiceName
	}

	price := sub.Price
	if input.Price != nil {
		price = *input.Price
	}

	userID := sub.UserID
	if input.UserID != nil {
		userID = *input.UserID
	}

	startDate := sub.StartDate
	if input.StartDate != nil {
		startDate = *input.StartDate
	}

	var endDate *string
	if input.EndDate != nil {
		endDate = input.EndDate
	} else {
		endDate = sub.EndDate
	}

	var endDateSQL sql.NullString
	if endDate != nil {
		endDateSQL = sql.NullString{String: *endDate, Valid: true}
	}

	err = r.db.QueryRow(ctx, query,
		serviceName, price, userID, startDate, endDateSQL, time.Now(), id,
	).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("failed to update subscription", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	r.logger.Info("subscription updated", zap.String("id", id))
	return sub, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM subscriptions WHERE id = $1"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete subscription", zap.Error(err), zap.String("id", id))
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrSubscriptionNotFound
	}

	r.logger.Info("subscription deleted", zap.String("id", id))
	return nil
}

// GetTotalCostForPeriod считает сумму цен подписок, активных хотя бы один день в периоде [startDate, endDate].
// Подписка активна в периоде, если: start_date <= period_end AND (end_date IS NULL OR end_date >= period_start).
func (r *SubscriptionRepository) GetTotalCostForPeriod(ctx context.Context, userID, serviceName, startDate, endDate string) (*models.TotalCostResponse, error) {
	query := `
		SELECT COALESCE(SUM(price), 0), COUNT(*)
		FROM subscriptions
		WHERE ($1::text IS NULL OR user_id = $1)
		  AND ($2::text IS NULL OR service_name = $2)
		  AND start_date <= $4
		  AND (end_date IS NULL OR end_date >= $3)
	`

	var totalCost int64
	var count int

	err := r.db.QueryRow(ctx, query,
		nullIfEmpty(userID),
		nullIfEmpty(serviceName),
		startDate,
		endDate,
	).Scan(&totalCost, &count)

	if err != nil {
		r.logger.Error("failed to calculate total cost", zap.Error(err))
		return nil, err
	}

	r.logger.Info("calculated total cost",
		zap.String("user_id", userID),
		zap.String("service_name", serviceName),
		zap.Int64("total_cost", totalCost),
	)

	return &models.TotalCostResponse{
		TotalCost: int(totalCost),
		Count:     count,
	}, nil
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
