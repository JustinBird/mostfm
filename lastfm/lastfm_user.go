package lastfm

import (
	"fmt"
	"net/http"
	"encoding/xml"
	"io"
)

func GetRecentTracks(secrets Secrets, user string) (LastFMRecentTracks, error) {
	var rt LastFMRecentTracks
	fields := []Field {
		{"api_key", secrets.APIKey},
		{"method", "user.getrecenttracks"},
		{"user", user},
	}

	resp, err := http.Get(createURL(fields))
	if err != nil {
		fmt.Println("Failed to get session!\n")
		return rt, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body!\n")
		return rt, err
	}

	err = xml.Unmarshal(body, &rt)
	if err != nil {
		fmt.Println("Failed to parse response!\n")
		return rt, err
	}

	if  rt.Status != "ok" {
		fmt.Printf("Bad status when getting token: %s\n", rt.Status)
	}

	return rt, nil
}