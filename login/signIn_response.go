package login

type SignInResponse struct {
	Error        bool            `json:"error"`
	AuthRequired bool            `json:"authenRequired,omitempty"`
	Message      string          `json:"message"`
	SignInUrl    string          `json:"signInUrl,omitempty"`
	Session      SessionResponse `json:"session,omitempty"`
}
