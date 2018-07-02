package authenticator

import "github.com/linkernetworks/web-services/entity"

type Mongo struct {
}

func NewMongo(url string) (*Mongo, error) {
	return &Mongo{}, nil
}

func (*Mongo) Verify(id, secret string) bool {
	// TODO: implement
	return true
}

func (m *Mongo) Register(user entity.User) error {
	// TODO: implement
	return nil
}
