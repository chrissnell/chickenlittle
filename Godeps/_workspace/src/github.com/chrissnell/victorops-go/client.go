package victorops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type API struct {
	APIKey   string
	Endpoint string
}

type Event struct {
	RoutingKey            string      `json:"-"`
	MessageType           MessageType `json:"message_type"`
	EntityID              string      `json:"entity_id,omitempty"`
	Timestamp             time.Time   `json:"-"`
	VOTimestamp           uint32      `json:"timestamp,omitempty"`
	StateStartTimestamp   time.Time   `json:"-"`
	VOStateStartTimestamp uint32      `json:"state_start_time,omitempty"`
	StateMessage          string      `json:"state_message,omitempty"`
	EntityIsHost          bool        `json:"entity_is_host,omitempty"`
	MonitoringTool        string      `json:"monitoring_tool,omitempty"`
	EntityDisplayName     string      `json:"entity_display_name,omitempty"`
	AckMsg                string      `json:"ack_msg,omitempty"`
	AckAuthor             string      `json:"ack_author,omitempty"`
}

type Response struct {
	Result   string `json:"result"`
	EntityID string `json:"entity_id"`
	Message  string `json:"message"`
}

type (
	MessageType string
)

const (
	Info            MessageType = "INFO"
	Warning         MessageType = "WARNING"
	Acknowledgement MessageType = "ACKNOWLEDGEMENT"
	Critical        MessageType = "CRITICAL"
	Recovery        MessageType = "RECOVERY"
)

func NewClient(apikey string) *API {
	a := &API{
		APIKey:   apikey,
		Endpoint: "https://alert.victorops.com/integrations/generic/20131114/alert",
	}
	return a
}

func (a *API) CreateEvent(rk string) *Event {
	e := &Event{
		RoutingKey: rk,
	}
	return e
}

func (a *API) SendAlert(e *Event) (*Response, error) {
	var r Response

	if !e.Timestamp.IsZero() {
		e.VOTimestamp = uint32(e.Timestamp.Unix())
	}
	if !e.StateStartTimestamp.IsZero() {
		e.VOStateStartTimestamp = uint32(e.StateStartTimestamp.Unix())
	}

	endpoint := fmt.Sprint(a.Endpoint, "/", a.APIKey, "/", e.RoutingKey)

	js, _ := json.Marshal(e)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(js))
	if err != nil {
		fmt.Println("Error:", err)
	}

	if resp.StatusCode != 200 {
		r := &Response{
			Result:  "failure",
			Message: resp.Status,
		}
		return r, fmt.Errorf("Error.  Verify your API key and the endpoint URL: %v", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &r)

	return &r, nil

}
