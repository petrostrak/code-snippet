package main

import (
	"html"
	"regexp"
	"testing"
)

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='{{.CSRFToken}}'>`)

func extractCSRFToken(t *testing.T, body []byte) string {
	// Use the FindSubmatch method to extract the token from the HTML body.
	// Note that this returns an array with the entire matched pattern in the
	// first position, and the values of any captured data in the subsequent
	// positions.
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}
