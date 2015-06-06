package main

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/ckvist/twilio/twiml"
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

func GenerateTwiML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if _, exists := NIP.Stoppers[uuid]; !exists {
		http.Error(w, "No active notifications for this UUID", http.StatusNotFound)
		return
	}

	resp := twiml.NewResponse()

	intro := twiml.Say{
		Voice: "woman",
		Text:  "This is Chicken Little with a message for you.",
	}

	gather := twiml.Gather{
		Action:    "http://me/uuid/digits",
		Timeout:   15,
		NumDigits: 1,
	}

	theMessage := twiml.Say{
		Voice: "man",
		Text:  NIP.Messages[uuid],
	}

	pressAny := twiml.Say{
		Voice: "woman",
		Text:  "Press any key to acknowledge receipt of this message",
	}

	resp.Action(intro)
	resp.Gather(gather, theMessage, pressAny)
	resp.Send(w)
}
