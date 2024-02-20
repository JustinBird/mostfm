package lastfm

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCreateURL(t *testing.T) {
	fields := []Field{
		{"api_key", "0123456789"},
		{"method", "test"},
		{"token", "abcdef"},
	}

	result := createURL(fields)
	expected := "http://ws.audioscrobbler.com/2.0/?api_key=0123456789&method=test&token=abcdef&"
	if result != expected {
		t.Errorf("Creating URL FAILED!\nExpected: %s\nActual: %s", expected, result)
	} else {
		t.Logf("Creating URL PASSED!")
	}
}

func TestCreateSignature(t *testing.T) {
	fields := []Field{
		{"api_key", "0123456789"},
		{"method", "test"},
		{"token", "abcdef"},
	}

	shared_secret := "password"

	createSignature(&fields, shared_secret)
	for i, field := range fields {
		if field.Key == "api_sig" {
			expected := "4fd2ca924382ed08de8e96d1bdf80de8"
			if field.Value != expected {
				t.Errorf("Creating signature FAILED!\nExpected: %s\nActual: %s", expected, field.Value)
			} else {
				t.Logf("Creating signature PASSED!")
			}
			break
		} else if i == len(fields) {
			t.Errorf("Creating signature FAILED!\nNo \"api_sig\" member!")
		}
	}
}

func TestNewAPIFromFile(t *testing.T) {
	tests := []struct {
		FileContents   string
		ExpectedAPIKey string
		ExpectedSecret string
		ExpectedError  error
	}{
		{
			FileContents:   "",
			ExpectedAPIKey: "",
			ExpectedSecret: "",
			ExpectedError:  ErrFileRead,
		},
		{
			FileContents:   "<secrets><apikey>API KEY GOES HERE</apikey><secret>SHARED SECRET GOES HERE</secret></secrets>",
			ExpectedAPIKey: "API KEY GOES HERE",
			ExpectedSecret: "SHARED SECRET GOES HERE",
			ExpectedError:  nil,
		},
		{
			FileContents:   "<secrets><apikey>API KEY GOES HERE<secret>SHARED SECRET GOES HERE</secrets>",
			ExpectedAPIKey: "",
			ExpectedSecret: "",
			ExpectedError:  ErrXMLParse,
		},
	}

	for i, test := range tests {
		filename := "does_not_exist.xml"
		if len(test.FileContents) != 0 {
			f, err := os.CreateTemp("", "apifile")
			if err != nil {
				t.Errorf("Failed to create temp file! %s", err)
			}
			defer os.Remove(f.Name())
			filename = f.Name()

			_, err = f.Write([]byte(test.FileContents))
			if err != nil {
				t.Errorf("Failed to write to temp file! %s", err)
			}
		}

		s, err := NewAPIFromFile(filename)
		if !errors.Is(err, test.ExpectedError) {
			t.Errorf("Test #%d: Expected error message did not match!\nExpected: %s\nActual: %s", i, test.ExpectedError, err)
		}

		if s.APIKey != test.ExpectedAPIKey || s.Secret != test.ExpectedSecret {
			t.Errorf("Test #%d: API key or secret is wrong!\nExpected: (%s, %s)\nActual: (%s, %s)", i, test.ExpectedAPIKey, test.ExpectedSecret, s.APIKey, s.Secret)
		}
	}
}

var api = LastFMAPI{APIKey: "0123456789", Secret: "abcdefg"}

func TestGetToken(t *testing.T) {
	tests := []struct {
		API           LastFMAPI
		Output        string
		HTTPCode      int
		ExpectedToken LastFMToken
		ExpectedError error
	}{
		{ // Normal operation
			API:      api,
			Output:   `<?xml version="1.0" encoding="UTF-8"?><lfm status="ok"><token>thisismytoken</token></lfm>`,
			HTTPCode: http.StatusOK,
			ExpectedToken: LastFMToken{
				XMLName: xml.Name{Local: "lfm"},
				Status:  "ok",
				Token:   "thisismytoken",
				Error: LastFMError{
					ErrorMsg:  "",
					ErrorCode: 0,
				},
			},
			ExpectedError: nil,
		},
		{ // Bad HTTP Status
			API:           api,
			Output:        "",
			HTTPCode:      http.StatusBadRequest,
			ExpectedToken: LastFMToken{},
			ExpectedError: ErrHTTPCode,
		},
		{ // Bad XML Parse
			API:      api,
			Output:   `<?xml version="1.0" encoding="UTF-8"?><lfm status="ok"><token>thisismytoken</token>`,
			HTTPCode: http.StatusOK,
			ExpectedToken: LastFMToken{
				XMLName: xml.Name{Local: "lfm"},
				Status:  "ok",
				Token:   "thisismytoken",
				Error: LastFMError{
					ErrorMsg:  "",
					ErrorCode: 0,
				},
			},
			ExpectedError: ErrXMLParse,
		},
		{ // Bad Last.fm status
			API:      api,
			Output:   `<?xml version="1.0" encoding="UTF-8"?><lfm status="fail"><token>thisismytoken</token></lfm>`,
			HTTPCode: http.StatusOK,
			ExpectedToken: LastFMToken{
				XMLName: xml.Name{Local: "lfm"},
				Status:  "fail",
				Token:   "thisismytoken",
				Error: LastFMError{
					ErrorMsg:  "",
					ErrorCode: 0,
				},
			},
			ExpectedError: ErrLastFMStatus,
		},
	}

	for i, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.HTTPCode)
			w.Write([]byte(test.Output))
		}))

		LastFMURL = server.URL
		token, err := test.API.GetToken()
		if !errors.Is(err, test.ExpectedError) {
			t.Errorf("Test #%d: Expected error message did not match!\nExpected: %s\nActual: %s", i, test.ExpectedError, err)
		}

		if !reflect.DeepEqual(token, test.ExpectedToken) {
			t.Errorf("Test #%d: LastFM token did not match expected!", i)
		}
	}
}

func TestGetSession(t *testing.T) {
	tests := []struct {
		API             LastFMAPI
		Token           string
		Output          string
		HTTPCode        int
		ExpectedSession LastFMSession
		ExpectedError   error
	}{
		{ // Normal operation
			API:      api,
			Token:    "token",
			Output:   `<?xml version="1.0" encoding="UTF-8"?><lfm status="ok"><session><name>mostfm</name><key>session-key</key><subscriber>1</subscriber></session></lfm>`,
			HTTPCode: http.StatusOK,
			ExpectedSession: LastFMSession{
				XMLName:    xml.Name{Local: "lfm"},
				Status:     "ok",
				Name:       "mostfm",
				Key:        "session-key",
				Subscriber: 1,
				Error: LastFMError{
					ErrorMsg:  "",
					ErrorCode: 0,
				},
			},
			ExpectedError: nil,
		},
		{ // Bad HTTP Status
			API:             api,
			Token:           "token",
			Output:          "",
			HTTPCode:        http.StatusBadRequest,
			ExpectedSession: LastFMSession{},
			ExpectedError:   ErrHTTPCode,
		},
		{ // Bad HTTP Status but output
			API:      api,
			Token:    "token",
			Output:   `<?xml version="1.0" encoding="UTF-8"?><lfm status="ok"><session><name>mostfm</name><key>session-key</key><subscriber>1</subscriber></session></lfm>`,
			HTTPCode: http.StatusBadRequest,
			ExpectedSession: LastFMSession{
				XMLName:    xml.Name{Local: "lfm"},
				Status:     "ok",
				Name:       "mostfm",
				Key:        "session-key",
				Subscriber: 1,
				Error: LastFMError{
					ErrorMsg:  "",
					ErrorCode: 0,
				},
			},
			ExpectedError: ErrHTTPCode,
		},
		{ // Bad XML Parse
			API:      api,
			Token:    "token",
			Output:   `<?xml version="1.0" encoding="UTF-8"?><lfm status="ok"><session><name>mostfm</name><key>session-key</key><subscriber>0</subscriber></session>`,
			HTTPCode: http.StatusOK,
			ExpectedSession: LastFMSession{
				XMLName:    xml.Name{Local: "lfm"},
				Status:     "ok",
				Name:       "mostfm",
				Key:        "session-key",
				Subscriber: 0,
				Error: LastFMError{
					ErrorMsg:  "",
					ErrorCode: 0,
				},
			},
			ExpectedError: ErrXMLParse,
		},
		{ // Bad Last.fm status with error
			API:      api,
			Token:    "token",
			Output:   `<?xml version="1.0" encoding="UTF-8"?><lfm status="failed"><error code="14">Unauthorized Token - This token has not been authorized</error></lfm>`,
			HTTPCode: http.StatusOK,
			ExpectedSession: LastFMSession{
				XMLName:    xml.Name{Local: "lfm"},
				Status:     "failed",
				Name:       "",
				Key:        "",
				Subscriber: 0,
				Error: LastFMError{
					XMLName:   xml.Name{Local: "error"},
					ErrorMsg:  "Unauthorized Token - This token has not been authorized",
					ErrorCode: 14,
				},
			},
			ExpectedError: ErrLastFMStatus,
		},
	}

	for i, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.HTTPCode)
			w.Write([]byte(test.Output))
		}))

		LastFMURL = server.URL
		session, err := test.API.GetSession(test.Token)
		if !errors.Is(err, test.ExpectedError) {
			t.Errorf("Test #%d: Expected error message did not match!\nExpected: %s\nActual: %s", i, test.ExpectedError, err)
		}

		if !reflect.DeepEqual(session, test.ExpectedSession) {
			t.Errorf("Test #%d: LastFM session did not match expected!", i)
		}
	}
}
