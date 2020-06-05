package response

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Content string `json:"message"`
}

var internalServerErrorMessage = `{"error":"internal server error"}`

func Error(w http.ResponseWriter, e error, code int) {
	w.WriteHeader(code)
	if e == nil {
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&Message{Content: e.Error()}); err != nil {
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
	}
}

func JSON(w http.ResponseWriter, body interface{}, code int) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		Error(w, err, http.StatusInternalServerError)
	}
}
