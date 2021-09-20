package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/zebresel-com/mongodm"
)

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Print("Send JSON err: ", err)
	}
}

func ErrorResponse(w http.ResponseWriter, status int, err ...error) {
	var errText string
	for _, e := range err {
		errText += e.Error() + ", "
	}
	errText = strings.TrimRight(errText, ", ")
	data := map[string]interface{}{
		"success": false,
		"error":   errText,
		"data":    nil,
	}
	SendJSON(w, status, data)
}

func ErrorHandler(w http.ResponseWriter, err error) {
	if _, ok := err.(*mongodm.NotFoundError); ok {
		ErrorResponse(w, http.StatusBadRequest, errors.New("not found resource"))
	} else if v, ok := err.(*mongodm.ValidationError); ok {
		ErrorResponse(w, http.StatusBadRequest, v)
	} else if v, ok := err.(*mongodm.DuplicateError); ok {
		ErrorResponse(w, http.StatusBadRequest, v)
	} else {
		ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
	}
}
