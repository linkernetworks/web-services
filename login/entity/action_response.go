package entity

import "github.com/linkernetworks/validator"

type ActionResponse struct {
	Error       bool                    `json:"error"`
	Validations validator.ValidationMap `json:"validations,omitempty"`
	Message     string                  `json:"message"`
}
