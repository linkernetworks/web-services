package userstorage

import (
	"testing"

	"github.com/linkernetworks/webservice/login/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyMemoryStorage(t *testing.T) {

	// arrange
	s := NewMemoryStorage()

	// action & assert
	assert.Nil(t, s.FindByEmail("aa"))
	assert.Nil(t, s.FindByToken("aa"))
	assert.Nil(t, s.FindByPassword("aa", "aa"))
}

func TestMemoryStorage(t *testing.T) {

	// arrange
	s := NewMemoryStorage()
	expected := &entity.User{
		Email:       "user@example.com",
		AccessToken: "tokennnn",
		Password:    "passssss",
	}
	err := s.Save(expected)
	require.NoError(t, err)

	// action & assert
	assert.Equal(t, expected, s.FindByEmail("user@example.com"))
	assert.Equal(t, expected, s.FindByToken("tokennnn"))
	assert.Equal(t, expected, s.FindByPassword("user@example.com", "passssss"))
}
