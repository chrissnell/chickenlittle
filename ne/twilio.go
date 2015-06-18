package ne

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
)

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

// Sends an SMS text message to a phone number using the Twilio API,
// optionally including a method for acknowledging receipt of the message.
func (e *Engine) SendSMS(phoneNumber, message, uuid string, dontSendAckRequest bool) {
	var cr SMSResponse

	if uuid != "" {
		log.Println("[", uuid, "]", "Sending SMS to", phoneNumber, "with message:", message)
	} else {
		log.Println("Sending SMS to", phoneNumber, "with message:", message)
	}

	// Generate an int in the range 100 <= n <= 999
	ackReply := rand.Intn(899) + 100

	// Builds a form that will be posted to Twilio API
	u := url.Values{}
	u.Set("From", e.Config.Twilio.CallFromNumber)
	u.Set("To", phoneNumber)

	// Sometimes we send texts that don't require ACKing.  This handles that.
	if dontSendAckRequest {
		u.Set("Body", message)
	} else {
		u.Set("Body", fmt.Sprint(message, " - Reply with \"", ackReply, "\" to acknowledge"))

	}

	// If we have a UUID, we can request status callbacks for this SMS
	if uuid != "" {
		u.Set("StatusCallback", fmt.Sprint(e.Config.Service.CallbackURLBase, "/", uuid, "/callback"))
	}

	// Post the request to the Twilio API
	body := *strings.NewReader(u.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprint(e.Config.Twilio.APIBaseURL, e.Config.Twilio.AccountSID, "/Messages.json"), &body)
	req.SetBasicAuth(e.Config.Twilio.AccountSID, e.Config.Twilio.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("SendSMS() Request error:", err)
	}

	// Get the response
	b, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024*20))
	resp.Body.Close()

	err = json.Unmarshal(b, &cr)
	if err != nil {
		log.Fatalln("SendSMS() Error unmarshalling JSON:", err)
	}

	if uuid != "" {
		// We create conversation key that's a combination of our recipient's phone number and the random 3-digit key
		// that we generated above
		conversationKey := fmt.Sprint(cr.To, "::", ackReply)

		e.SetConversation(conversationKey, uuid)
	}

}

// Makes a phone call to a phone number using the Twilio API.  Sends Twilio a URL for
// retrieving the TwiML that defines the interaction in the call.
func (e *Engine) MakePhoneCall(phoneNumber, message, uuid string) {
	var cr map[string]interface{}

	log.Println("[", uuid, "] Calling", phoneNumber, "with message:", message)

	// Build a form that we'll POST to the Twilio API to initiate a phone call
	u := url.Values{}
	u.Set("From", e.Config.Twilio.CallFromNumber)
	u.Set("To", phoneNumber)
	u.Set("Url", fmt.Sprint(e.Config.Service.CallbackURLBase, "/", uuid, "/twiml/notify"))
	// Optional status callbacks are enabled below...
	// u.Set("StatusCallback", fmt.Sprint(c.Config.Service.CallbackURLBase, "/", uuid, "/callback"))
	// u.Add("StatusCallbackEvent", "ringing")
	// u.Add("StatusCallbackEvent", "answered")
	// u.Add("StatusCallbackEvent", "completed")
	u.Set("IfMachine", "Hangup")
	u.Set("Timeout", "20")
	body := *strings.NewReader(u.Encode())

	// Send our form to Twilio
	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprint(e.Config.Twilio.APIBaseURL, e.Config.Twilio.AccountSID, "/Calls.json"), &body)
	req.SetBasicAuth(e.Config.Twilio.AccountSID, e.Config.Twilio.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("MakePhoneCall() Request error:", err)
	}

	b, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024*20))
	resp.Body.Close()

	// We get the response back but don't currently do anything with it.   TO DO: implement error handling
	err = json.Unmarshal(b, &cr)
	if err != nil {
		log.Fatalln("MakePhoneCall() Error unmarshalling JSON:", err)
	}

}
