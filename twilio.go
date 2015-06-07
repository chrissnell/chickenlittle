package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"bitbucket.org/ckvist/twilio/twiml"
	"github.com/gorilla/mux"
)

type CallbackResponse struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func MakePhoneCall(phoneNumber, message, uuid string) {
	var cr map[string]interface{}

	log.Println("Calling", phoneNumber, "with message:", message)

	u := url.Values{}
	u.Set("From", c.Config.Integrations.Twilio.CallFromNumber)
	u.Set("To", phoneNumber)
	u.Set("Url", fmt.Sprint(c.Config.Service.CallbackURLBase, "/", uuid, "/twiml"))
	u.Set("StatusCallback", fmt.Sprint(c.Config.Service.CallbackURLBase, "/", uuid, "/callback"))
	u.Add("StatusCallbackEvent", "ringing")
	u.Add("StatusCallbackEvent", "answered")
	u.Add("StatusCallbackEvent", "completed")
	u.Set("IfMachine", "Hangup")
	u.Set("Timeout", "20")
	body := *strings.NewReader(u.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprint(c.Config.Integrations.Twilio.APIBaseURL, c.Config.Integrations.Twilio.AccountSID, "/Calls.json"), &body)
	req.SetBasicAuth(c.Config.Integrations.Twilio.AccountSID, c.Config.Integrations.Twilio.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("MakePhoneCall() Request error:", err)
	}

	b, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024*20))
	resp.Body.Close()

	err = json.Unmarshal(b, &cr)
	if err != nil {
		log.Fatalln("MakePhoneCall() Error unmarshalling JSON:", err)
	}

	log.Printf("Call Response Received:\n%+v\n", cr)

}

func ReceiveCallback(w http.ResponseWriter, r *http.Request) {
	var res CallbackResponse

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Stuff will happen

	res = CallbackResponse{
		Message: "Callback received",
		UUID:    uuid,
	}

	json.NewEncoder(w).Encode(res)

}

func ReceiveDigits(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	err := r.ParseForm()
	if err != nil {
		log.Println("ReceiveDigits() r.ParseForm() error:", err)
	}

	digits := r.FormValue("Digits")

	log.Println("Digits received:", digits)

	if digits != "" {

		if _, exists := NIP.Stoppers[uuid]; !exists {
			log.Println("ReceiveDigits(): No active notifications for this UUID:", uuid)
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// Attempt to stop the notification by sending the UUID to the notification engine
		stopChan <- uuid
	}
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
		Action:    fmt.Sprint(c.Config.Service.CallbackURLBase, "/", uuid, "/digits"),
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
