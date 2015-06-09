package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
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

type SMSResponse struct {
	Sid         string
	DateCreated string
	DateUpdated string
	DateSent    string
	AccountSid  string
	To          string
	From        string
	Body        string
	NumSegments string
	Status      string
	Direction   string
	Price       string
	PriceUnit   string
	ApiVersion  string
	Uri         string
}

func SendSMS(phoneNumber, message, uuid string, dontSendAckRequest bool) {
	var cr SMSResponse

	if uuid != "" {
		log.Println("[", uuid, "]", "Sending SMS to", phoneNumber, "with message:", message)
	} else {
		log.Println("Sending SMS to", phoneNumber, "with message:", message)
	}

	// Generate an int in the range 100 <= n <= 999
	ackReply := rand.Intn(899) + 100

	u := url.Values{}
	u.Set("From", c.Config.Integrations.Twilio.CallFromNumber)
	u.Set("To", phoneNumber)
	if dontSendAckRequest {
		u.Set("Body", message)
	} else {
		u.Set("Body", fmt.Sprint(message, " - Reply with \"", ackReply, "\" to acknowledge"))

	}

	if uuid != "" {
		u.Set("StatusCallback", fmt.Sprint(c.Config.Service.CallbackURLBase, "/", uuid, "/callback"))
	}

	body := *strings.NewReader(u.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprint(c.Config.Integrations.Twilio.APIBaseURL, c.Config.Integrations.Twilio.AccountSID, "/Messages.json"), &body)
	req.SetBasicAuth(c.Config.Integrations.Twilio.AccountSID, c.Config.Integrations.Twilio.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("SendSMS() Request error:", err)
	}

	b, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024*20))
	resp.Body.Close()

	err = json.Unmarshal(b, &cr)
	if err != nil {
		log.Fatalln("SendSMS() Error unmarshalling JSON:", err)
	}

	if uuid != "" {
		// We make a conversation key that's a combination of our recipient's phone number and the random 3-digit key
		// that we generated abovee
		conversationKey := fmt.Sprint(cr.To, "::", ackReply)

		NIP.Mu.Lock()
		defer NIP.Mu.Unlock()

		NIP.Conversations[conversationKey] = uuid
	}

}

func MakePhoneCall(phoneNumber, message, uuid string) {
	var cr map[string]interface{}

	log.Println("[", uuid, "] Calling", phoneNumber, "with message:", message)

	u := url.Values{}
	u.Set("From", c.Config.Integrations.Twilio.CallFromNumber)
	u.Set("To", phoneNumber)
	u.Set("Url", fmt.Sprint(c.Config.Service.CallbackURLBase, "/", uuid, "/twiml/notify"))
	// u.Set("StatusCallback", fmt.Sprint(c.Config.Service.CallbackURLBase, "/", uuid, "/callback"))
	// u.Add("StatusCallbackEvent", "ringing")
	// u.Add("StatusCallbackEvent", "answered")
	// u.Add("StatusCallbackEvent", "completed")
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

}

func ReceiveSMSReply(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("ReceiveSMSReply() r.ParseForm() error:", err)
	}

	recipient := r.FormValue("From")
	if recipient == "" {
		log.Println("ReceiveSMSReply() error: 'From' parameter was not provided in response")
		return
	}

	NIP.Mu.Lock()

	conversationKey := fmt.Sprint(recipient, "::", r.FormValue("Body"))

	if _, exists := NIP.Conversations[conversationKey]; exists {
		uuid := NIP.Conversations[conversationKey]

		log.Println("[", uuid, "]", "Recieved a SMS reply from", recipient, ":", r.FormValue("Body"))

		if _, exists := NIP.Stoppers[uuid]; !exists {
			log.Println("ReceiveSMSReply(): No active notifications for this UUID:", uuid)
			http.Error(w, "", http.StatusNotFound)
			NIP.Mu.Unlock()
			return
		}

		// Delete the conversation key from the in-progress store
		delete(NIP.Conversations, conversationKey)

		// Unlock our mutex so the notification engine can take it
		NIP.Mu.Unlock()

		log.Println("[", uuid, "] Attempting to stop notifications")

		// Attempt to stop the notification by sending the UUID to the notification engine
		stopChan <- uuid

		SendSMS(recipient, "Chicken Little has received your acknowledgment.  Thanks!", uuid, true)

	} else {
		SendSMS(recipient, "I'm sorry but I don't recognize that response.   Please acknowledge with the three-digit code from the notfication you received.", "", true)
	}

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

	// Fetch some form values we'll need from Twilio's request
	digits := r.FormValue("Digits")
	callSid := r.FormValue("CallSid")

	// If digits has been set, user has answered the phone and pressed (any) key to acknowledge the message
	if digits != "" {

		if _, exists := NIP.Stoppers[uuid]; !exists {
			log.Println("ReceiveDigits(): No active notifications for this UUID:", uuid)
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// We matched a valid notification-in-progress and the user pressed digits when prompted
		// so we'll do a POST to Twilio that points the call at a TwiML routine that confirms
		// their acknowledgement and sends them on their way.
		u := url.Values{}
		u.Set("Url", fmt.Sprint(c.Config.Service.CallbackURLBase, "/", uuid, "/twiml/acknowledged"))

		body := *strings.NewReader(u.Encode())

		client := &http.Client{}
		req, _ := http.NewRequest("POST", fmt.Sprint(c.Config.Integrations.Twilio.APIBaseURL, c.Config.Integrations.Twilio.AccountSID, "/Calls/", callSid), &body)
		req.SetBasicAuth(c.Config.Integrations.Twilio.AccountSID, c.Config.Integrations.Twilio.AuthToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Send the POST request
		_, err := client.Do(req)
		if err != nil {
			log.Println("ReceiveDigits() TwiML POST Request error:", err)
		}

		// Attempt to stop the notification by sending the UUID to the notification engine
		stopChan <- uuid
	}
}

func GenerateTwiML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	action := vars["action"]

	resp := twiml.NewResponse()

	switch action {
	case "notify":
		// This is a request for a TwiML script for a standard message notification
		if _, exists := NIP.Stoppers[uuid]; !exists {
			http.Error(w, "No active notifications for this UUID", http.StatusNotFound)
			return
		}

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

	case "acknowledged":
		// This is a request for the end-of-call wrap-up message
		resp.Action(twiml.Say{
			Voice: "woman",
			Text:  "Thank you. This message has been acknowledged. Goodbye!",
		})
	}

	resp.Send(w)
}
