package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/linkernetworks/web-services/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignIn(t *testing.T) {
	// arrange: prepare HTTP request
	req, _ := http.NewRequest("POST", "/test_path", nil)
	w := httptest.NewRecorder()
	web, err := New(&config.Config{})

	// arrange: create web instance
	require.NotNil(t, web)
	require.NoError(t, err)

	// action
	web.SignIn(w, req)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}
