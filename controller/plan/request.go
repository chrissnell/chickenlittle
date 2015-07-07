package plan

import "github.com/chrissnell/chickenlittle/model"

// NotificationPlanResponse is a struct to encode an NotificationPlan as JSON
type NotificationPlanResponse struct {
	NotificationPlan model.NotificationPlan `json:"people"`
	Message          string                 `json:"message"`
	Error            string                 `json:"error"`
}
