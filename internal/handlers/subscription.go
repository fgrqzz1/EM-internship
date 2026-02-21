package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"em-internship/internal/models"
	"em-internship/internal/service"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
	logger  *zap.Logger
}

func NewSubscriptionHandler(service *service.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

// CreateSubscription godoc
// @Summary Create a subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.CreateSubscriptionInput true "Subscription data"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var input models.CreateSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	sub, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err))
		http.Error(w, `{"error":"failed to create subscription"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Get subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Warn("subscription not found", zap.String("id", id), zap.Error(err))
		http.Error(w, `{"error":"subscription not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(sub)
}

// ListSubscriptions godoc
// @Summary List all subscriptions
// @Description Get paginated list of subscriptions
// @Tags subscriptions
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.SubscriptionList
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	list, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to get subscriptions", zap.Error(err))
		http.Error(w, `{"error":"failed to get subscriptions"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(list)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update an existing subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body models.UpdateSubscriptionInput true "Subscription data"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input models.UpdateSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	sub, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		h.logger.Warn("failed to update subscription", zap.String("id", id), zap.Error(err))
		http.Error(w, `{"error":"failed to update subscription"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(sub)
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Warn("failed to delete subscription", zap.String("id", id), zap.Error(err))
		http.Error(w, `{"error":"subscription not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTotalCost godoc
// @Summary Get total cost for period
// @Description Calculate total cost of subscriptions for a period with filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID filter"
// @Param service_name query string false "Service name filter"
// @Param start_date query string true "Start date (MM-YYYY)"
// @Param end_date query string true "End date (MM-YYYY)"
// @Success 200 {object} models.TotalCostResponse
// @Router /subscriptions/total-cost [get]
func (h *SubscriptionHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		http.Error(w, `{"error":"start_date and end_date are required"}`, http.StatusBadRequest)
		return
	}

	response, err := h.service.GetTotalCostForPeriod(r.Context(), userID, serviceName, startDate, endDate)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			http.Error(w, `{"error":"invalid date format: use MM-YYYY"}`, http.StatusBadRequest)
			return
		}
		h.logger.Error("failed to calculate total cost", zap.Error(err))
		http.Error(w, `{"error":"failed to calculate total cost"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}
