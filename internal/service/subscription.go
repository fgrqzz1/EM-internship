package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"em-internship/internal/models"
	"em-internship/internal/repository"
)

type SubscriptionService struct {
	repo      *repository.SubscriptionRepository
	validator *validator.Validate
	logger    *zap.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, logger *zap.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:      repo,
		validator: validator.New(),
		logger:    logger,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, input models.CreateSubscriptionInput) (*models.Subscription, error) {
	if err := s.validator.StructCtx(ctx, input); err != nil {
		s.logger.Warn("validation error", zap.Error(err))
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return s.repo.Create(ctx, input)
}

func (s *SubscriptionService) GetByID(ctx context.Context, id string) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubscriptionService) GetAll(ctx context.Context, limit, offset int) (*models.SubscriptionList, error) {
	if limit <= 0 { // простая проверка на дурака
		limit = 15
	} else if limit > 99 {
		limit = 99
	}

	return s.repo.GetAll(ctx, limit, offset)
}

func (s *SubscriptionService) Update(ctx context.Context, id string, input models.UpdateSubscriptionInput) (*models.Subscription, error) {
	return s.repo.Update(ctx, id, input)
}

func (s *SubscriptionService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubscriptionService) GetTotalCostForPeriod(ctx context.Context, userID, serviceName, startDate, endDate string) (*models.TotalCostResponse, error) {
	return s.repo.GetTotalCostForPeriod(ctx, userID, serviceName, startDate, endDate)
}
