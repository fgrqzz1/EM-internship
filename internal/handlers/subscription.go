package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"

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

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Warn("failed to delete subscription", zap.String("id", id), zap.Error(err))
		http.Error(w, `{"error":"subscription not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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
		h.logger.Error("failed to calculate total cost", zap.Error(err))
		http.Error(w, `{"error":"failed to calculate total cost"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}
