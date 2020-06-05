package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iseroukhov/brave-new-billing/pkg/entities/payment"
	"github.com/iseroukhov/brave-new-billing/pkg/http/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PaymentRepositoryInterface interface {
	Create(amount float32, purpose string) (*payment.Payment, error)
	GetByUID(uid string) (*payment.Payment, error)
	SetStatus(p *payment.Payment, statusID int64) error
	All() ([]*payment.Payment, error)
}

type RegisterHandler struct {
	logger *logrus.Logger
	repo   *payment.Repository //PaymentRepositoryInterface
}

func NewRegisterHandler(logger *logrus.Logger, repo *payment.Repository) *RegisterHandler {
	return &RegisterHandler{
		logger: logger,
		repo:   repo,
	}
}

type RegisterResponse struct {
	UID       uuid.UUID  `json:"uid"`
	FormURL   string     `json:"formUrl"`
	ErrorInfo *ErrorInfo `json:"errorInfo"`
}

type ErrorInfo struct {
	Message string `json:"message"`
}

func (h *RegisterHandler) Index() http.HandlerFunc {
	type requestBody struct {
		Amount         float32 `json:"amount"`
		PaymentPurpose string  `json:"payment_purpose"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body := &requestBody{}
		registerResponse := &RegisterResponse{
			ErrorInfo: &ErrorInfo{},
		}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			registerResponse.ErrorInfo.Message = "internal server error"
			response.JSON(w, registerResponse, http.StatusInternalServerError)
			h.logger.Errorf("register.Index (Decode): %s", err.Error())
			return
		}

		if body.Amount <= 0.0 {
			registerResponse.ErrorInfo.Message = "amount field is required"
			response.JSON(w, registerResponse, http.StatusBadRequest)
			return
		}

		if body.PaymentPurpose == "" {
			registerResponse.ErrorInfo.Message = "payment_purpose field is required"
			response.JSON(w, registerResponse, http.StatusBadRequest)
			return
		}

		p, err := h.repo.Create(body.Amount, body.PaymentPurpose)
		if err != nil {
			registerResponse.ErrorInfo.Message = "internal server error"
			response.JSON(w, registerResponse, http.StatusInternalServerError)
			h.logger.Errorf("register.Index (Index): %s", err.Error())
			return
		}

		registerResponse.UID = p.UID
		registerResponse.FormURL = fmt.Sprintf("http://localhost:9000/payments/card/form?sessionId=%s", p.UID)
		registerResponse.ErrorInfo = nil

		response.JSON(w, registerResponse, http.StatusOK)
	}
}
