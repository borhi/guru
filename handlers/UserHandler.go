package handlers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"guru/models"
	"guru/services"
	"net/http"
)

type UserHandler struct {
	service   *services.UserService
	validator *validator.Validate
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
		validator: validator.New(),
	}
}

func (h *UserHandler) Get(w http.ResponseWriter, req *http.Request) {
	var userRequest models.GetUserRequestModel
	var errorResponse models.ErrorResponseModel
	err := json.NewDecoder(req.Body).Decode(&userRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Error = err.Error()
		zap.L().Error(err.Error())
	}

	err = h.validator.Struct(&userRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Error = err.Error()
		zap.L().Error(err.Error())
	}

	userResponse, err := h.service.GetUser(userRequest.Id, userRequest.Token)
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
		}
		if err.Error() == "wrong token" {
			w.WriteHeader(http.StatusBadRequest)
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

	if err := json.NewEncoder(w).Encode(userResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zap.L().Error(err.Error())
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, req *http.Request) {
	var userRequest models.UserModel
	var res models.ErrorResponseModel
	err := json.NewDecoder(req.Body).Decode(&userRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = err.Error()
		zap.L().Error(err.Error())
	}

	err = h.validator.Struct(&userRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Error = err.Error()
		zap.L().Error(err.Error())
	}

	user := models.UserModel{Id: userRequest.Id, Balance: userRequest.Balance, Token: userRequest.Token}
	h.service.CreateUser(userRequest.Id, user)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zap.L().Error(err.Error())
	}
}

func (h *UserHandler) AddDeposit(w http.ResponseWriter, req *http.Request) {
	var depositRequest models.DepositRequestModel
	var errorResponse models.ErrorResponseModel

	err := json.NewDecoder(req.Body).Decode(&depositRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Error = err.Error()
		zap.L().Error(err.Error())
	}

	err = h.validator.Struct(&depositRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse.Error = err.Error()
		zap.L().Error(err.Error())
	}

	depositResponse, err := h.service.AddDeposit(depositRequest)
	if err != nil {
		switch err.Error() {
		case "wrong token":
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

	if err := json.NewEncoder(w).Encode(depositResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zap.L().Error(err.Error())
	}
}
