package api

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	w := httptest.NewRecorder()
	Ping(w, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	result, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, "\"Pong\"", string(result))
}
