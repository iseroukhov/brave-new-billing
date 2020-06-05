package handlers

import (
	"encoding/json"
	"github.com/iseroukhov/brave-new-billing/pkg/entities/payment"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PaymentHandler struct {
	logger *logrus.Logger
	repo   *payment.Repository
}

func NewPaymentHandler(logger *logrus.Logger, repo *payment.Repository) *PaymentHandler {
	return &PaymentHandler{
		logger: logger,
		repo:   repo,
	}
}

func (h *PaymentHandler) Index() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		payments, err := h.repo.All(r.FormValue("from"), r.FormValue("to"))
		if err != nil {
			h.internalServerError(w, "payment.Index (All)", err)
			return
		}

		err = json.NewEncoder(w).Encode(payments)
		if err != nil {
			h.internalServerError(w, "payment.Index (Encode)", err)
			return
		}
	}
}

func (h *PaymentHandler) internalServerError(w http.ResponseWriter, message string, err error) {
	h.logger.Error(message+":", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	err = templates.ExecuteTemplate(w, "500.html", nil)
	if err != nil {
		h.logger.Error("internal server error: " + err.Error())
		http.Error(w, `internal server error`, http.StatusInternalServerError)
	}
}
