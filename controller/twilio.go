package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"bitbucket.org/ckvist/twilio/twiml"
	"github.com/gorilla/mux"
)

// CallbackResponse is the json struct as returned from twilio in the callbacks
type CallbackResponse struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// ReceiveSMSReply receives the SMS reply callback from Twilio and deletes the notification if the
// response text matches the code sent with the original SMS notification
func (a *Controller) ReceiveSMSReply(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("ReceiveSMSReply() r.ParseForm() error:", err)
	}

	// We should have a "From" parameter being passed from Twilio
	recipient := r.FormValue("From")
	if recipient == "" {
		log.Println("ReceiveSMSReply() error: 'From' parameter was not provided in response")
		return
	}

	// Our conversation key is a combination of the recipient's phone number and the 3-digit code
	// that they sent in reply
	conversationKey := fmt.Sprint(recipient, "::", r.FormValue("Body"))

	uuid, err := a.n.GetConversation(conversationKey)
	if err != nil {
		a.n.SendSMS(recipient, "I'm sorry but I don't recognize that response.   Please acknowledge with the three-digit code from the notfication you received.", "", true)
	}
	log.Println("[", uuid, "]", "Recieved a SMS reply from", recipient, ":", r.FormValue("Body"))

	if !a.n.IsNotification(uuid) {
		log.Println("ReceiveSMSReply(): No active notifications for this UUID:", uuid)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	a.n.RemoveConversation(conversationKey)
	log.Println("[", uuid, "] Attempting to stop notifications")
	a.n.CancelNotification(uuid)
	a.n.SendSMS(recipient, "Chicken Little has received your acknowledgment.  Thanks!", uuid, true)

}

// ReceiveCallback receives call progress callbacks from the Twilio API.  Not currently used.
// May be used for Websocket interface in the future.
func (a *Controller) ReceiveCallback(w http.ResponseWriter, r *http.Request) {
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

// ReceiveDigits receives digits pressed during a phone call via callback by the Twilio API.
// Stops the notification if the user pressed any keys.
func (a *Controller) ReceiveDigits(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	err := r.ParseForm()
	if err != nil {
		log.Println("ReceiveDigits() r.ParseForm() error:", err)
	}

	// Fetch some form values we'll need from Twilio's request
	digits := r.FormValue("Digits")
	callSid := r.FormValue("CallSid")

	// If digits has been set, user has answered the phone and pressed (any) key to acknowledge the message
	if digits != "" {

		if !a.n.IsNotification(uuid) {
			log.Println("ReceiveDigits(): No active notifications for this UUID:", uuid)
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// We matched a valid notification-in-progress and the user pressed digits when prompted
		// so we'll do a POST to Twilio that points the call at a TwiML routine that confirms
		// their acknowledgement and sends them on their way.
		u := url.Values{}
		u.Set("Url", fmt.Sprint(a.c.Service.CallbackURLBase, "/", uuid, "/twiml/acknowledged"))

		// Send our POST to Twilio
		body := *strings.NewReader(u.Encode())
		client := &http.Client{}
		req, _ := http.NewRequest("POST", fmt.Sprint(a.c.Integrations.Twilio.APIBaseURL, a.c.Integrations.Twilio.AccountSID, "/Calls/", callSid), &body)
		req.SetBasicAuth(a.c.Integrations.Twilio.AccountSID, a.c.Integrations.Twilio.AuthToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Send the POST request
		_, err := client.Do(req)
		if err != nil {
			log.Println("ReceiveDigits() TwiML POST Request error:", err)
		}

		// Attempt to stop the notification by sending the UUID to the notification engine
		a.n.CancelNotification(uuid)
	}
}

// GenerateTwiML is a Twilio callback which generates TwiML that is used to describe the flow of the phone call.
func (a *Controller) GenerateTwiML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	action := vars["action"]

	resp := twiml.NewResponse()

	switch action {
	case "notify":
		// This is a request for a TwiML script for a standard message notification
		if !a.n.IsNotification(uuid) {
			http.Error(w, "No active notifications for this UUID", http.StatusNotFound)
			return
		}

		intro := twiml.Say{
			Voice: "woman",
			Text:  "This is Chicken Little with a message for you.",
		}

		gather := twiml.Gather{
			Action:    fmt.Sprint(a.c.Service.CallbackURLBase, "/", uuid, "/digits"),
			Timeout:   15,
			NumDigits: 1,
		}

		theMessage := twiml.Say{
			Voice: "man",
			Text:  a.n.GetMessage(uuid),
		}

		pressAny := twiml.Say{
			Voice: "woman",
			Text:  "Press any key to acknowledge receipt of this message",
		}

		resp.Action(intro)
		resp.Gather(gather, theMessage, pressAny)

	case "acknowledged":
		// This is a request for the end-of-call wrap-up message
		resp.Action(twiml.Say{
			Voice: "woman",
			Text:  "Thank you. This message has been acknowledged. Goodbye!",
		})
	}

	// Reply to the callback with the TwiML content
	resp.Send(w)
}
