package login

import "github.com/linkernetworks/webservice/login/entity"

type SessionResponse struct {
	ID          string      `json:"id,omitempty"`
	Token       string      `json:"token,omitempty"`
	ExpiredAt   int64       `json:"expiredAt,omitempty"`
	CurrentUser entity.User `json:"currentUser,omitempty"`
}
