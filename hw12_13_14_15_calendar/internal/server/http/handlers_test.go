package internalhttp

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Hello(t *testing.T) {
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	req, err := http.NewRequest("GET", "/hello", nil)
	require.NoError(t, err)
	mux.ServeHTTP(rr, req)

	require.Equal(t, "Hello world!", rr.Body.String())
	require.Equal(t, 200, rr.Code)
}
