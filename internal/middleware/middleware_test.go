package middleware

import (
	"fmt"
	"net/http"
	"testing"
)

func TestContentTypeMiddleware(t *testing.T) {
	var testHandler TestHandler
	mH := ContentTypeMiddleware(&testHandler)

	switch v := mH.(type) {
	case http.Handler:
		// Do nothing as this is the type we expect
	default:
		t.Error(fmt.Sprintf("%T received unexpected type for ContentTypeMiddleware", v))
	}
}
