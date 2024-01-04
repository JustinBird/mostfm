package lastfm

import (
	"fmt"
	"net/http"
	"encoding/xml"
	"io"
	"os"
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
	resp, err := http.Get(createURL(fields))
	if err != nil {
		fmt.Println("Failed to get token!")
		return t, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return t, err
	}

	err = xml.Unmarshal(body, &t)
	if err != nil {
		return t, err
	}

	if  t.Status != "ok" {
		fmt.Printf("Bad status when getting token: %s\n", t.Status)
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

	resp, err := http.Get(createURL(fields))
	if err != nil {
		fmt.Println("Failed to get session!")
		return s, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	err = xml.Unmarshal(body, &s)
	if err != nil {
		return s, err
	}

	if  s.Status != "ok" {
		fmt.Printf("Bad status when getting token: %s\n", s.Status)
	}

	return s, nil
}