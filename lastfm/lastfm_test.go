package lastfm


import "testing"

func TestCreateURL(t *testing.T) {
	fields := []Field {
		{"api_key", "0123456789"},
		{"method", "test"},
		{"token", "abcdef"},
	}

	result := createURL(fields)
	expected := "http://ws.audioscrobbler.com/2.0/?api_key=0123456789&method=test&token=abcdef&"
	if result == expected {
		t.Errorf("Creating URL FAILED!\nExpected: %s\nActual: %s", expected, result)
	} else {
		t.Logf("Creating URL PASSED!")
	}
}