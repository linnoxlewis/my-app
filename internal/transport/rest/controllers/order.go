package controllers

import (
	"applicationDesignTest/internal/models"
	"applicationDesignTest/internal/models/dto"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

var errInvalidJson = errors.New("invalid json format")

type response struct {
	Success bool          `json:"success"`
	Data    *models.Order `json:"data"`
	Error   string        `json:"error"`
}

type orderService interface {
	CreateOrder(ctx context.Context, order dto.Order) (*models.Order, error)
}

type OrderController struct {
	orderSrv orderService
}

func NewOrderController(orderSrv orderService) *OrderController {
	return &OrderController{orderSrv}
}

func (o *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder dto.Order
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{
			Success: false,
			Error:   errInvalidJson.Error(),
		})
		return
	}

	if err := newOrder.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	order, err := o.orderSrv.CreateOrder(context.Background(), newOrder)
	if err != nil {
		msgErr, status := o.getError(err)
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response{
			Success: false,
			Error:   msgErr,
		})
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response{
		Success: true,
		Data:    order,
	})
}

func (o *OrderController) getError(err error) (error string, status int) {
	if errors.Is(err, models.ErrNotAvailableRooms) {
		return err.Error(), http.StatusConflict
	} else if errors.Is(err, models.ErrNotFoundInformation) {
		return err.Error(), http.StatusNotFound
	} else {
		return models.ErrInternalServerError.Error(), http.StatusInternalServerError
	}
}
