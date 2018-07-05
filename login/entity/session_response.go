package entity

type SessionResponse struct {
	ID          string `json:"id,omitempty"`
	Token       string `json:"token,omitempty"`
	ExpiredAt   int64  `json:"expiredAt,omitempty"`
	CurrentUser User   `json:"currentUser,omitempty"`
}
