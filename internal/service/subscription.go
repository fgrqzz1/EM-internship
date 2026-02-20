package service

import (
	"context"
	"em-internship/internal/models"
	"em-internship/internal/repository"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type SubscriptionService struct {
	repo      *repository.SubscriptionRepository
	validator *validator.Validator
	logger    *zap.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, logger *zap.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:      repo,
		validator: validator.New(),
		logger:    logger,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, input models.CreateSubscriptionInput) (*models.Subscription, error)) {
	if err := s.validator.StructCtx(ctx, input); err != nil {
		s.logger.Warn("validation error", zap.Error(err))
		return nil, fmt.Error("validation error", err)
	}

	return s.repo.Create(ctx, input)
}

func
