package lastfm

import (
	"fmt"
	"encoding/xml"
	"os"
	"errors"
)

func GetSecrets(secrets_path string) (Secrets, error) {
	var s Secrets
	data, err := os.ReadFile(secrets_path)
	if err != nil {
		return s, err
	}

	err = xml.Unmarshal(data, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func GetToken(apikey string) (LastFMToken, error) {
	var t LastFMToken
	fields := []Field {
		{"api_key", apikey},
		{"method",  "auth.getToken"},
	}

	err := LastFMCall(&fields, &t)
	if err != nil {
		return t, fmt.Errorf("Failed to get token: %w", err)
	}

	if  t.Status != "ok" {
		return t, fmt.Errorf("%w Status: %s", ErrLastFMStatus, t.Status)
	}

	return t, nil
}

func GetSession(secrets Secrets, token string) (LastFMSession, error) {
	var s LastFMSession
	fields := []Field {
		{"api_key", secrets.APIKey},
		{"method", "auth.getSession"},
		{"token", token},
	}
	createSignature(&fields, secrets.Secret)

	err := LastFMCall(&fields, &s)
	if err != nil {
		err := errors.Join(err, errors.New("Failed to get session!"))
		return s, err
	}

	if  s.Status != "ok" {
		fmt.Printf("Bad status when getting session: %s\n", s.Status)
	}

	return s, nil
}