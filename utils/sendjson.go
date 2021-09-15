package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
