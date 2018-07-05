package entity

import oauth "github.com/linkernetworks/oauth/entity"

type SessionResponse struct {
	ID          string     `json:"id,omitempty"`
	Token       string     `json:"token,omitempty"`
	ExpiredAt   int64      `json:"expiredAt,omitempty"`
	CurrentUser oauth.User `json:"currentUser,omitempty"`
}
