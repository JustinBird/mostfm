package lastfm

import (
	"fmt"
	//"net/http"
	//"encoding/xml"
	//"io"
	"errors"
)

func (api LastFMAPI) GetRecentTracks(user string) (LastFMRecentTracks, error) {
	var rt LastFMRecentTracks
	fields := []Field{
		{"api_key", api.APIKey},
		{"method", "user.getrecenttracks"},
		{"user", user},
	}

	err := LastFMCall(&fields, &rt)
	if err != nil {
		err := errors.Join(err, errors.New("failed to get recent tracks"))
		return rt, err
	}

	if rt.Status != "ok" {
		fmt.Printf("Bad status when getting recent tracks: %s\n", rt.Status)
	}

	return rt, nil
}
