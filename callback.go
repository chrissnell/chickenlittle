package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CallbackResponse struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ReceiveCallback(w http.ResponseWriter, r *http.Request) {
	var res CallbackResponse

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	// Stuff will happen

	res = CallbackResponse{
		Message: "Callback received",
		UUID:    uuid,
	}

	json.NewEncoder(w).Encode(res)

}
