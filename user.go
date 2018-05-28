package acl

import (
	"fmt"
	"net/http"

	oauth "github.com/linkernetworks/oauth/entity"
	"bitbucket.org/linkernetworks/aurora/src/service/session"
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/mongo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const SessionKey = "ses"

func GetCurrentUserRestful(ses *mongo.Session, req *restful.Request) (*oauth.User, error) {
	token := req.Request.Header.Get("Authorization")
	if len(token) == 0 {
		return GetCurrentUser(ses, req.Request)
	}

	return GetCurrentUserByToken(ses, token)
}

// GetCurrentUser get current user data with login session and return user data
// excluding sensitive data like password.
func GetCurrentUser(ses *mongo.Session, req *http.Request) (*oauth.User, error) {
	email, err := GetCurrentUserEmail(req)
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
func GetCurrentUserByToken(ses *mongo.Session, token string) (*oauth.User, error) {
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
func GetCurrentUserWithPassword(ses *mongo.Session, req *http.Request) (*oauth.User, error) {
	email, err := GetCurrentUserEmail(req)
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

func GetCurrentUserEmail(req *http.Request) (string, error) {
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
