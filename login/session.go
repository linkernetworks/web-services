package login

import (
	"errors"
	"fmt"
	"net/http"

	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/mongo"
	response "github.com/linkernetworks/net/http"
	oauth "github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/session"
	"github.com/linkernetworks/webservice/login/entity"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

func (s *LoginService) registerUserSession(req *http.Request, resp http.ResponseWriter, ses *sessions.Session, u *oauth.User) error {

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

func (s *LoginService) signInSession(req *http.Request, resp http.ResponseWriter, user *oauth.User) (uuid.UUID, error) {
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

func (s *LoginService) GetCurrentUserRestful(ses *mongo.Session, req *restful.Request) (*oauth.User, error) {
	token := req.Request.Header.Get("Authorization")
	if len(token) == 0 {
		return s.GetCurrentUser(ses, req.Request)
	}

	return s.GetCurrentUserByToken(ses, token)
}

// GetCurrentUser get current user data with login session and return user data
// excluding sensitive data like password.
func (s *LoginService) GetCurrentUser(ses *mongo.Session, req *http.Request) (*oauth.User, error) {
	email, err := s.GetCurrentUserEmail(req)
	if err != nil {
		return nil, err
	}

	user := oauth.User{}
	q := bson.M{"email": email}
	projection := bson.M{"password": 0}
	if err := ses.C(oauth.UserCollectionName).Find(q).Select(projection).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, fmt.Errorf("user document not found.")
		}
		return nil, err
	}

	return &user, nil
}

// GetCurrentUserByToken get current user data with login token and return user data
// excluding sensitive data like password.
func (s *LoginService) GetCurrentUserByToken(ses *mongo.Session, token string) (*oauth.User, error) {
	user := oauth.User{}
	q := bson.M{"access_token": token}
	projection := bson.M{"password": 0}
	if err := ses.C(oauth.UserCollectionName).Find(q).Select(projection).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, fmt.Errorf("user document not found.")
		}
		return nil, err
	}

	return &user, nil
}

// GetCurrentUserWithPassword get current user data with login session and return all user data
// including sensitive data like encrypted password.
func (s *LoginService) GetCurrentUserWithPassword(ses *mongo.Session, req *http.Request) (*oauth.User, error) {
	email, err := s.GetCurrentUserEmail(req)
	if err != nil {
		return nil, err
	}

	user := oauth.User{}
	q := bson.M{"email": email}
	if err := ses.C(oauth.UserCollectionName).Find(q).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, fmt.Errorf("user document not found.")
		}
		return nil, err
	}

	return &user, nil
}

func (s *LoginService) GetCurrentUserEmail(req *http.Request) (string, error) {
	session, err := session.Service.Store.Get(req, SessionKey)
	if err != nil {
		return "", err
	}

	val, found := session.Values["email"]
	if !found {
		return "", fmt.Errorf("session email is not set.")
	}

	email, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("session email value type is invalid.")
	}
	return email, err
}
