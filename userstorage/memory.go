package userstorage

import (
	"sync"

	"github.com/linkernetworks/webservice/login/entity"
)

type MemoryStorage struct {
	lock    sync.Mutex
	byMail  map[string]entity.User
	byToken map[string]entity.User
	byPass  map[string]entity.User // key: "user@example.cpm_password"
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		byMail:  make(map[string]entity.User),
		byToken: make(map[string]entity.User),
		byPass:  make(map[string]entity.User),
	}
}

func (s *MemoryStorage) FindByEmail(email string) *entity.User {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, ok := s.byMail[email]
	if !ok {
		return nil
	}
	return &user
}

func (s *MemoryStorage) FindByToken(token string) *entity.User {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, ok := s.byToken[token]
	if !ok {
		return nil
	}
	return &user
}

func (s *MemoryStorage) FindByPassword(email, password string) *entity.User {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, ok := s.byPass[email+"_"+password]
	if !ok {
		return nil
	}
	return &user
}

func (s *MemoryStorage) Save(user *entity.User) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if user.Email != "" {
		s.byMail[user.Email] = *user
	}

	if user.AccessToken != "" {
		s.byToken[user.AccessToken] = *user
	}

	if user.Email != "" && user.Password != "" {
		s.byPass[user.Email+"_"+user.Password] = *user
	}

	return nil
}
