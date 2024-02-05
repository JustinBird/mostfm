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
		fmt.Println("Failed to get session!")
		return rt, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	//fmt.Printf(string(body))
	xml.Unmarshal(body, &rt)

	return rt, nil
}