package service

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"

	"em-internship/internal/repository"
)

func TestGetTotalCostForPeriod_InvalidDateFormat(t *testing.T) {
	logger := zap.NewNop()
	svc := NewSubscriptionService((*repository.SubscriptionRepository)(nil), logger)
	ctx := context.Background()

	tests := []struct {
		name      string
		startDate string
		endDate   string
	}{
		{"invalid start_date month", "13-2025", "12-2025"},
		{"invalid end_date format", "01-2025", "2025-01"},
		{"empty start_date", "", "12-2025"},
		{"empty end_date", "01-2025", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetTotalCostForPeriod(ctx, "", "", tt.startDate, tt.endDate)
			if err == nil {
				t.Fatal("expected error")
			}
			if !errors.Is(err, ErrInvalidDateFormat) {
				t.Errorf("expected ErrInvalidDateFormat, got %v", err)
			}
		})
	}
}
