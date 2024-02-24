package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMsg struct {
	Msg     string `json:"error"`
	Success bool   `json:"success"`
}

func ResError(err error, w http.ResponseWriter, code int, message string) {
	log.Println(err)
	ResJSON(w, code, ErrorMsg{Msg: message, Success: false})
}

func ResJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
