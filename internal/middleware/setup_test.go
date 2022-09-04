package middleware

import (
	"net/http"
	"os"
	"testing"
)

// This runs before every test
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type TestHandler struct{}

func (th *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
