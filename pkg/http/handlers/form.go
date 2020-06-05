package handlers

import (
	"github.com/iseroukhov/brave-new-billing/pkg/entities/payment"
	"github.com/iseroukhov/brave-new-billing/pkg/luhn"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type FormHandler struct {
	logger *logrus.Logger
	repo   *payment.Repository
}

func NewFormHandler(logger *logrus.Logger, repo *payment.Repository) *FormHandler {
	return &FormHandler{
		logger: logger,
		repo:   repo,
	}
}

var templates = template.Must(template.ParseFiles(
	"./templates/form.html",
	"./templates/success.html",
	"./templates/error.html",
	"./templates/404.html",
	"./templates/500.html"),
)

func (h *FormHandler) Index() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.FormValue("sessionId")
		p, err := h.repo.GetByUID(sessionId)
		if err != nil {
			h.internalServerError(w, "form.Index (GetByUID)", err)
			return
		}
		if p == nil {
			h.notFound(w)
			return
		}
		err = templates.ExecuteTemplate(w, "form.html", p)
		if err != nil {
			h.internalServerError(w, "form.Index (ExecuteTemplate)", err)
		}
	}
}

func (h *FormHandler) Store() http.HandlerFunc {
	type responseType struct {
		Header string
		Text   string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.FormValue("sessionId")
		p, err := h.repo.GetByUID(sessionId)
		if err != nil {
			h.internalServerError(w, "form.Store (GetByUID)", err)
			return
		}
		if p == nil {
			h.notFound(w)
			return
		}
		if luhn.IsValid(r.FormValue("card_number")) {
			if err := h.repo.SetStatus(p, payment.StatusPaid); err != nil {
				h.internalServerError(w, "form valid (SetStatus)", err)
				return
			}
			err = templates.ExecuteTemplate(w, "success.html", &responseType{
				Text: "Оплата прошла успешно",
			})
			if err != nil {
				h.internalServerError(w, "form.Store (ExecuteTemplate)", err)
			}
			return
		}

		if err := h.repo.SetStatus(p, payment.StatusError); err != nil {
			h.internalServerError(w, "form not valid (SetStatus)", err)
			return
		}
		err = templates.ExecuteTemplate(w, "error.html", &responseType{
			Header: "Ошибка при оплате",
			Text:   "Неверный номер карты",
		})
		if err != nil {
			h.internalServerError(w, "form.Store (ExecuteTemplate)", err)
			return
		}
	}
}

func (h *FormHandler) notFound(w http.ResponseWriter) {
	err := templates.ExecuteTemplate(w, "404.html", nil)
	if err != nil {
		h.logger.Error("form.notFound (internal server error): " + err.Error())
		http.Error(w, `internal server error`, http.StatusInternalServerError)
	}
}

func (h *FormHandler) internalServerError(w http.ResponseWriter, message string, err error) {
	h.logger.Error(message+":", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	err = templates.ExecuteTemplate(w, "500.html", nil)
	if err != nil {
		h.logger.Error("internal server error: " + err.Error())
		http.Error(w, `internal server error`, http.StatusInternalServerError)
	}
}
