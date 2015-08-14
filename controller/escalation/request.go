package escalation

import "github.com/chrissnell/chickenlittle/model"

// Response is used to marshal the responses to JSON
type Response struct {
	Plans   []model.EscalationPlan `json:"plans"`
	Message string                 `json:"message"`
	Error   string                 `json:"error"`
}
