package rotation

import "github.com/chrissnell/chickenlittle/model"

// Response is used to marshal the responses
type Response struct {
	Policies []model.RotationPolicy `json:"policies"`
	Message  string                 `json:"message"`
	Error    string                 `json:"error"`
}
