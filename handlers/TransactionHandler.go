package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"guru/models"
	"guru/services"
	"net/http"
)

type TransactionHandler struct {
	service   *services.UserService
	validator *validator.Validate
}

func NewTransactionHandler(service *services.UserService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		validator: validator.New(),
	}
}

func (h *TransactionHandler) Transaction(w http.ResponseWriter, req *http.Request) {
	var transactionRequest models.TransactionRequestModel
	var errorResponse models.ErrorResponseModel

	err := json.NewDecoder(req.Body).Decode(&transactionRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Error = err.Error()
		zap.L().Error(err.Error())
	}

	err = h.validator.Struct(&transactionRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Error = err.Error()
		zap.L().Error(err.Error())
	}

	transactionResponse, err := h.service.Transaction(transactionRequest)
	if err != nil {
		switch err.Error() {
		case "wrong token", "not enough balance":
			w.WriteHeader(http.StatusBadRequest)
		case "not found":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		errorResponse.Error = err.Error()
		zap.L().Error(err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	if errorResponse.Error != "" {
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			zap.L().Error(err.Error())
		}
		return
	}

	if err := json.NewEncoder(w).Encode(transactionResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zap.L().Error(err.Error())
	}
}
