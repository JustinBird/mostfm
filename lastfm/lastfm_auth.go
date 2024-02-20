package lastfm

import (
	"errors"
	"fmt"
)

func (api LastFMAPI) GetToken() (LastFMToken, error) {
	var t LastFMToken
	fields := []Field{
		{"api_key", api.APIKey},
		{"method", "auth.getToken"},
	}

	err := LastFMCall(&fields, &t)
	if err != nil {
		return t, fmt.Errorf("failed to get token: %w", err)
	}

	if t.Status != "ok" {
		return t, fmt.Errorf("%w Status: %s", ErrLastFMStatus, t.Status)
	}

	return t, nil
}

func (api LastFMAPI) GetSession(token string) (LastFMSession, error) {
	var s LastFMSession
	fields := []Field{
		{"api_key", api.APIKey},
		{"method", "auth.getSession"},
		{"token", token},
	}
	createSignature(&fields, api.Secret)

	err := LastFMCall(&fields, &s)
	if err != nil {
		err := errors.Join(err, errors.New("failed to get session"))
		return s, err
	}

	if s.Status != "ok" {
		return s, fmt.Errorf("%w Status: %s", ErrLastFMStatus, s.Status)
	}

	return s, nil
}
