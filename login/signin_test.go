package login

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/webservice/login/entity"
	"github.com/linkernetworks/webservice/pwdutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SignInTestSuite struct {
	suite.Suite
	ls *LoginService
	wc *restful.Container
}

func (s *SignInTestSuite) SetupTest() {
	var err error

	tmpDir, err := ioutil.TempDir("", "SignInTestSuite")
	require.NoError(s.T(), err)

	logger.Setup(logger.LoggerConfig{
		Dir:   tmpDir + "/log/",
		Level: "debug",
	})

	s.ls, err = New(nil)
	require.NoError(s.T(), err)

	s.wc = restful.NewContainer()
	s.wc.Add(s.ls.WebService())
}

func TestSignInTestSuite(t *testing.T) {
	suite.Run(t, new(SignInTestSuite))
}

// As a valide user, I can sign-in successfully.
func (s *SignInTestSuite) TestSignInWithValidUser() {
	// arrange: create a dummy user in storage
	encPass, err := pwdutil.EncryptPasswordLegacy("aaaaaa", s.ls.passworldSalt)
	require.NoError(s.T(), err)
	s.ls.userStorage.Save(&entity.User{
		Password: encPass,
		Email:    "user@example.com",
		Roles:    []string{"simple_role"},
	})

	// arrange: prepare HTTP request
	req, _ := http.NewRequest("POST", "/signin", strings.NewReader(`
		 {
			"password": "aaaaaa",
			"email": "user@example.com"
		 }
	`))
	req.Header.Add("Content-Type", "application/json")

	// action
	w := httptest.NewRecorder()
	s.wc.ServeHTTP(w, req)

	// assert
	assert.Equal(s.T(), http.StatusOK, w.Code)
	// TODO: check other fields in JSON responce
}
