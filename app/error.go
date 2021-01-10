package app

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewMessage(msg string, code int, resp http.ResponseWriter) (http.ResponseWriter, error) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	err := json.NewEncoder(resp).Encode(Message{
		Message: msg,
		Code:    code,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func BindResponse(js interface{}, resp http.ResponseWriter, status int) (http.ResponseWriter, error) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header()
	resp.WriteHeader(status)
	if js != nil {
		err := json.NewEncoder(resp).Encode(js)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}
