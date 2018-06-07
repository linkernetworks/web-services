package session

import (
	"errors"
	"net/http"

	"time"

	"bitbucket.org/linkernetworks/aurora/src/validator"
	oauth "github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/session"

	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
)

// will be the cookie name defined in the http header
const SessionKey = "ses"

type ActionResponse struct {
	Error       bool                    `json:"error"`
	Validations validator.ValidationMap `json:"validations,omitempty"`
	Message     string                  `json:"message"`
}

type SignInResponse struct {
	Error        bool            `json:"error"`
	AuthRequired bool            `json:"authenRequired,omitempty"`
	Message      string          `json:"message"`
	SignInUrl    string          `json:"signInUrl,omitempty"`
	Session      SessionResponse `json:"session,omitempty"`
}

type SessionResponse struct {
	ID          string     `json:"id,omitempty"`
	Token       string     `json:"token,omitempty"`
	ExpiredAt   int64      `json:"expiredAt,omitempty"`
	CurrentUser oauth.User `json:"currentUser,omitempty"`
}

// func NewLoginService(sp *serviceprovider.Container) *restful.WebService {
// 	ws := new(restful.WebService)
// 	ws.Path("/v1").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
// 	ws.Route(ws.GET("/me").Filter(sessionAuthenticationFilter).To(RESTfulServiceHandler(sp, GetMeHandler)))
// 	ws.Route(ws.POST("/email/check").To(RESTfulServiceHandler(sp, CheckEmailAvailability)))
// 	ws.Route(ws.POST("/signup").To(RESTfulServiceHandler(sp, SignUpUserHandler)))
// 	ws.Route(ws.POST("/signin").To(RESTfulServiceHandler(sp, SignInUserHandler)))
// 	ws.Route(ws.GET("/signout").Filter(sessionAuthenticationFilter).To(RESTfulServiceHandler(sp, SignOutUserHandler)))
// 	return ws
// }

func AllocateNewSessionToken() uuid.UUID {
	return uuid.NewV4()
}

func GetSession(req *http.Request) (*sessions.Session, error) {
	return session.Service.Store.Get(req, SessionKey)
}

func RegisterUserSession(req *http.Request, resp http.ResponseWriter, ses *sessions.Session, u *oauth.User) error {

	if len(u.Email) == 0 {
		return errors.New("email is required to register user session.")
	}

	if len(u.Roles) == 0 {
		return errors.New("at least one role is required.")
	}

	ses.Values["email"] = u.Email
	ses.Values["roles"] = u.Roles
	ses.Values["expiredAt"] = time.Now().Add(time.Minute * time.Duration(60)).Unix()
	return ses.Save(req, resp)
}

func SignIn(req *http.Request, resp http.ResponseWriter, user *oauth.User) (uuid.UUID, error) {
	token := AllocateNewSessionToken()
	ses, err := GetSession(req)
	if err != nil {
		return token, err
	}
	return token, RegisterUserSession(req, resp, ses, user)
}
