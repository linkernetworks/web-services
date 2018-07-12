package userstorage

import (
	"github.com/linkernetworks/webservice/login/entity"
)

const ERROR_NOTEXIST string = "ERROR_NOTEXIST"

type UserStorage interface {
	FindByEmail(email string) *entity.User
	FindByToken(token string) *entity.User
	FindByPassword(email, password string) *entity.User
	Save(user *entity.User) error
}
