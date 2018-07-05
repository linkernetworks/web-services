package usersession

import (
	"errors"
	"fmt"
	"net/http"

	"time"

	restful "github.com/emicklei/go-restful"
	response "github.com/linkernetworks/net/http"
	oauth "github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/session"
	"github.com/linkernetworks/webservice/login/entity"

	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
)

// will be the cookie name defined in the http header
const SessionKey = "ses"

func AllocateNewSessionToken() uuid.UUID {
	return uuid.NewV4()
}

func GetSession(req *http.Request) (*sessions.Session, error) {
	return session.Service.Store.Get(req, SessionKey)
}

// Pre-Handler user session authentication
func SessionAuthenticationFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	sess, err := session.Service.Store.Get(req.Request, SessionKey)
	if err != nil {
		msg := fmt.Errorf("Redis get auth token failed: %v", err)
		response.InternalServerError(req.Request, resp.ResponseWriter, msg)
		return
	}

	if !isExpired(sess) {
		//refresh
		sess.Values["expiredAt"] = time.Now().Add(24 * time.Hour).Unix()
		if err := sess.Save(req.Request, resp); err != nil {
			response.InternalServerError(req.Request, resp.ResponseWriter, fmt.Errorf("Redis save auth token failed: %v", err))
			return
		}
		chain.ProcessFilter(req, resp)
		return
	}

	resp.WriteHeaderAndEntity(http.StatusForbidden, entity.SignInResponse{
		Error:     true,
		Message:   "Unauthorized. Redirect to signin page",
		SignInUrl: "/signin",
	})
	return
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

func isExpired(sess *sessions.Session) bool {
	expiredAt := sess.Values["expiredAt"]
	if expiredAt == nil {
		return true
	}
	return expiredAt.(int64) < time.Now().Unix()
}
