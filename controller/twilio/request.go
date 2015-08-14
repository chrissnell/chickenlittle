package twilio

// CallbackResponse is the json struct as returned from twilio in the callbacks
type CallbackResponse struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
