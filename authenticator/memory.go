package authenticator

import (
	"fmt"

	"github.com/linkernetworks/web-services/entity"
)

type Memory struct {
	users map[string]entity.User
}

func NewMemory() (*Memory, error) {
	return &Memory{
		users: make(map[string]entity.User),
	}, nil
}

func (m *Memory) Verify(id, secret string) bool {
	user, exist := m.users[id]
	if !exist {
		return false
	}
	return user.Password == secret
}

func (m *Memory) Register(user entity.User) error {
	if _, exist := m.users[user.Email]; exist {
		return fmt.Errorf("[%v] already exists", user.Email)
	}
	m.users[user.Email] = user
	return nil
}
