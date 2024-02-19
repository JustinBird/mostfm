package lastfm

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"errors"
	"reflect"
	"encoding/xml"
	"os"
)

func TestCreateURL(t *testing.T) {
	fields := []Field {
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
	fields := []Field {
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
	} {
		{
			FileContents: "",
			ExpectedAPIKey: "",
			ExpectedSecret: "",
			ExpectedError: ErrFileRead,
		},
		{
			FileContents: "<secrets><apikey>API KEY GOES HERE</apikey><secret>SHARED SECRET GOES HERE</secret></secrets>",
			ExpectedAPIKey: "API KEY GOES HERE",
			ExpectedSecret: "SHARED SECRET GOES HERE",
			ExpectedError: nil,
		},
		{
			FileContents: "<secrets><apikey>API KEY GOES HERE<secret>SHARED SECRET GOES HERE</secrets>",
			ExpectedAPIKey: "",
			ExpectedSecret: "",
			ExpectedError: ErrXMLParse,
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

func TestGetToken(t *testing.T) {
	tests := []struct {
		API                  LastFMAPI
		Output               string
		HTTPCode             int
		ExpectedStatus       string
		ExpectedToken        LastFMToken
		ExpectedLastFMError  LastFMError
		ExpectedError        error
	} {
		{ // Normal operation
			API: LastFMAPI{ APIKey: "0123456789", Secret: "abcdefg" },
			Output: `<?xml version="1.0" encoding="UTF-8"?><lfm status="ok"><token>thisismytoken</token></lfm>`,
			HTTPCode: http.StatusOK,
			ExpectedToken: LastFMToken{
				XMLName: xml.Name{ Local: "lfm" },
				Status: "ok",
				Token: "thisismytoken",
				Error: LastFMError{
					ErrorMsg: "",
					ErrorCode: 0,
				},
			},
			ExpectedError: nil,
		},
		{ // Bad HTTP Status
			API: LastFMAPI{ APIKey: "0123456789", Secret: "abcdefg" },
			Output: "",
			HTTPCode: http.StatusBadRequest,
			ExpectedToken: LastFMToken{},
			ExpectedError: ErrHTTPCode,
		},
		{ // Bad XML Parse
			API: LastFMAPI{ APIKey: "0123456789", Secret: "abcdefg" },
			Output: `<?xml version="1.0" encoding="UTF-8"?><lfm status="ok"><token>thisismytoken</token>`,
			HTTPCode: http.StatusOK,
			ExpectedToken: LastFMToken{
				XMLName: xml.Name{ Local: "lfm" },
				Status: "ok",
				Token: "thisismytoken",
				Error: LastFMError{
					ErrorMsg: "",
					ErrorCode: 0,
				},
			},
			ExpectedError: ErrXMLParse,
		},
		{ // Bad Last.fm status
			API: LastFMAPI{ APIKey: "0123456789", Secret: "abcdefg" },
			Output: `<?xml version="1.0" encoding="UTF-8"?><lfm status="fail"><token>thisismytoken</token></lfm>`,
			HTTPCode: http.StatusOK,
			ExpectedToken: LastFMToken{
				XMLName: xml.Name{ Local: "lfm" },
				Status: "fail",
				Token: "thisismytoken",
				Error: LastFMError{
					ErrorMsg: "",
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

		if (!reflect.DeepEqual(token, test.ExpectedToken)) {
			t.Errorf("Test #%d: LastFM token did not match expected!", i)
		}
	}
}