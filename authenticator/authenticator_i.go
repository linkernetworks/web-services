package authenticator

import "github.com/linkernetworks/web-services/entity"

type Authenticator interface {
	Verify(id, secret string) bool
	Register(user entity.User) error
}
