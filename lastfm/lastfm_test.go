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
	type GetSecretsTest struct {
		Filename             string
		ExpectedAPIKey       string
		ExpectedSecret       string
		ExpectedError        bool
		ExpectedErrorMessage string
	}
	
	tests := []GetSecretsTest {
		{
			Filename: "doesnt_exist.xml",
			ExpectedAPIKey: "",
			ExpectedSecret: "",
			ExpectedError: true,
			ExpectedErrorMessage: "open doesnt_exist.xml: no such file or directory",
		},
		{
			Filename: "../secrets_example.xml",
			ExpectedAPIKey: "API KEY GOES HERE",
			ExpectedSecret: "SHARED SECRET GOES HERE",
			ExpectedError: false,
			ExpectedErrorMessage: "",
		},
	}

	for i, test := range tests {
		s, err := GetSecrets(test.Filename)
		if test.ExpectedError && err == nil {
			t.Errorf("Test #%d: Expected an error but err was nil!", i)
		} else if test.ExpectedError && err.Error() != test.ExpectedErrorMessage {
			t.Errorf("Test #%d: Expected error message did not match!\nExpected: %s\nActual: %s", i, err.Error(), test.ExpectedErrorMessage)
		} else if !test.ExpectedError && err != nil {
			t.Errorf("Test #%d: Expected no error but err was not nil!", i)
		}
		
		if s.APIKey != test.ExpectedAPIKey || s.Secret != test.ExpectedSecret {
			t.Errorf("Test #%d: API key or secret is wrong!\nExpected: (%s, %s)\nActual: (%s, %s)", i, s.APIKey, s.Secret, test.ExpectedAPIKey, test.ExpectedSecret)
		} else {
			t.Logf("Test #%d: API key and secret is correct!", i)
		}
	}
}