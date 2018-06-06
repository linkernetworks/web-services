package session

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bitbucket.org/linkernetworks/aurora/src/aurora/serviceprovider"
	"bitbucket.org/linkernetworks/aurora/src/entity"
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/config"
	oauth "github.com/linkernetworks/oauth/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignInUserHandler(t *testing.T) {
	cf := config.MustRead("../config/testing.json")
	sp := serviceprovider.New(cf)
	msession := sp.Mongo.NewSession()
	defer msession.Close()

	user := createTestUser(t, msession)
	require.NotNil(t, user)
	defer msession.Remove(oauth.UserCollectionName, "_id", user.ID)

	form := oauth.User{
		Email:    user.Email,
		Password: "testtest",
	}
	defer msession.Remove(entity.DatasetCollectionName, "owner", user.ID)

	bodyBytes, err := json.MarshalIndent(form, "", "  ")
	assert.NoError(t, err)

	bodyReader := strings.NewReader(string(bodyBytes))
	httpRequest, err := http.NewRequest("POST", "http://here.com/v1/signin", bodyReader)
	httpRequest.Header.Add("Content-Type", "application/json")
	assert.NoError(t, err)

	httpWriter := httptest.NewRecorder()
	wc := restful.NewContainer()
	wc.Add(newLoginService(sp))
	wc.Dispatch(httpWriter, httpRequest)
	assertResponseCode(t, 200, httpWriter)
}
