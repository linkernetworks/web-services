package login

import (
	"errors"
	"fmt"
	"net/http"

	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	response "github.com/linkernetworks/net/http"
	"github.com/linkernetworks/session"
	"github.com/linkernetworks/webservice/login/entity"

	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
)

// will be the cookie name defined in the http header
const SessionKey = "ses"

func (s *LoginService) allocateNewSessionToken() uuid.UUID {
	return uuid.NewV4()
}

func (s *LoginService) getSession(req *http.Request) (*sessions.Session, error) {
	return s.store.Get(req, SessionKey)
}

// Pre-Handler user session authentication
func (s *LoginService) authenticatedFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	sess, err := s.store.Get(req.Request, SessionKey)
	if err != nil {
		msg := fmt.Errorf("Redis get auth token failed: %v", err)
		response.InternalServerError(req.Request, resp.ResponseWriter, msg)
		return
	}

	if !s.isExpired(sess) {
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

func (s *LoginService) registerUserSession(req *http.Request, resp http.ResponseWriter, ses *sessions.Session, u *entity.User) error {

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

func (s *LoginService) signInSession(req *http.Request, resp http.ResponseWriter, user *entity.User) (uuid.UUID, error) {
	token := s.allocateNewSessionToken()
	ses, err := s.getSession(req)
	if err != nil {
		return token, err
	}
	return token, s.registerUserSession(req, resp, ses, user)
}

func (s *LoginService) isExpired(sess *sessions.Session) bool {
	expiredAt := sess.Values["expiredAt"]
	if expiredAt == nil {
		return true
	}
	return expiredAt.(int64) < time.Now().Unix()
}

func (s *LoginService) GetCurrentUserRestful(req *restful.Request) *entity.User {
	token := req.Request.Header.Get("Authorization")
	if len(token) == 0 {
		return s.GetCurrentUser(req.Request)
	}

	user := s.userStorage.FindByToken(token)

	return user
}

// GetCurrentUser get current user data with login session and return user data
// excluding sensitive data like password.
func (s *LoginService) GetCurrentUser(req *http.Request) *entity.User {
	email := s.GetCurrentUserEmail(req)
	if email == "" {
		return nil
	}

	user := s.userStorage.FindByEmail(email)

	return user
}

func (s *LoginService) GetCurrentUserEmail(req *http.Request) string {
	session, err := session.Service.Store.Get(req, SessionKey)
	if err != nil {
		logger.Errorf("Get session [%v] failed. err: [%v]", SessionKey, err)
		return ""
	}

	val, found := session.Values["email"]
	if !found {
		return ""
	}

	email, ok := val.(string)
	if !ok {
		logger.Errorf("Convert [%v] to string failed.", val)
		return ""
	}
	return email
}
