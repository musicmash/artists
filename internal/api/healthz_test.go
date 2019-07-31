package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_Healthz_OK(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := http.Get(server.URL + "/healthz")
	defer func() { _ = resp.Body.Close() }()

	// arrange
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAPI_Healthz_Error(t *testing.T) {
	server = httptest.NewServer(getMux())
	defer server.Close()

	// action
	resp, err := http.Get(server.URL + "/healthz")
	defer func() { _ = resp.Body.Close() }()

	// arrange
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
