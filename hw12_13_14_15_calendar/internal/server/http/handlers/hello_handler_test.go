package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMyHandler(t *testing.T) {
	service := NewHelloService(nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", service.GetHello)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	cases := []struct {
		name         string
		method       string
		target       string
		body         io.Reader
		responseCode int
	}{
		{"bad_request", http.MethodPost, "/hello", nil, http.StatusMethodNotAllowed},
		{"ok", http.MethodGet, "/hello", nil, http.StatusOK},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, c.method, ts.URL+c.target, c.body)
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, c.responseCode, res.StatusCode)
		})
	}
}
