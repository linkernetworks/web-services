package acl

import (
	"net/http"

	"bitbucket.org/linkernetworks/aurora/src/entity"
	"bitbucket.org/linkernetworks/aurora/src/service/mongo"
	restful "github.com/emicklei/go-restful"
)

type Callback func(*mongo.Session, *http.Request) (*entity.User, error)

func GetCurrentUserRestful(getCurrentUser Callback, ses *mongo.Session, req *restful.Request) (*entity.User, error) {

	return getCurrentUser(ses, req.Request)
}
