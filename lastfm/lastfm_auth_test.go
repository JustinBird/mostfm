package lastfm

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"errors"
	"reflect"
	"encoding/xml"
	"io/fs"
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

func TestGetSecrets(t *testing.T) {
	tests := []struct {
		Filename             string
		ExpectedAPIKey       string
		ExpectedSecret       string
		ExpectedError        error
	} {
		{
			Filename: "doesnt_exist.xml",
			ExpectedAPIKey: "",
			ExpectedSecret: "",
			ExpectedError: &fs.PathError{},
		},
		{
			Filename: "../secrets_example.xml",
			ExpectedAPIKey: "API KEY GOES HERE",
			ExpectedSecret: "SHARED SECRET GOES HERE",
			ExpectedError: nil,
		},
		{
			Filename: "../secrets_broken.xml",
			ExpectedAPIKey: "",
			ExpectedSecret: "",
			ExpectedError: &xml.SyntaxError{},
		},
	}

	for i, test := range tests {
		t.Logf("%s", test.Filename)
		t.Logf("%s", reflect.TypeOf(test.ExpectedError))
		s, err := GetSecrets(test.Filename)
		t.Logf("%s", reflect.TypeOf(err))
		if reflect.TypeOf(err) != reflect.TypeOf(test.ExpectedError) {
			t.Errorf("Test #%d: Expected error did not match!\nExpected: %s\nActual: %s", i, reflect.TypeOf(test.ExpectedError),  reflect.TypeOf(err))
		}
		
		if s.APIKey != test.ExpectedAPIKey || s.Secret != test.ExpectedSecret {
			t.Errorf("Test #%d: API key or secret is wrong!\nExpected: (%s, %s)\nActual: (%s, %s)", i, test.ExpectedAPIKey, test.ExpectedSecret, s.APIKey, s.Secret)
		} else {
			t.Logf("Test #%d: API key and secret is correct!", i)
		}
	}
}

func TestGetToken(t *testing.T) {
	tests := []struct {
		APIKey               string
		Secret               string
		Output               string
		HTTPCode             int
		ExpectedStatus       string
		ExpectedToken        LastFMToken
		ExpectedLastFMError  LastFMError
		ExpectedError        error
	} {
		{ // Normal operation
			APIKey: "0123456789",
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
			APIKey: "0123456789",
			Output: "",
			HTTPCode: http.StatusBadRequest,
			ExpectedToken: LastFMToken{},
			ExpectedError: ErrHTTPCode,
		},
		{ // Bad XML Parse
			APIKey: "0123456789",
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
			APIKey: "0123456789",
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
		token, err := GetToken(test.APIKey)
		if !errors.Is(err, test.ExpectedError) {
			t.Errorf("Test #%d: Expected error message did not match!\nExpected: %s\nActual: %s", i, test.ExpectedError, err)
		}

		if (!reflect.DeepEqual(token, test.ExpectedToken)) {
			t.Errorf("Test #%d: LastFM token did not match expected!", i)
		}
	}
}