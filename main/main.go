package main

import (
	"bufio"
	"fmt"
	"os"

	"mostfm/lastfm"
)

func main() {
	api, err := lastfm.NewAPIFromFile("secrets.xml")
	if err != nil {
		fmt.Println("Failed to get secrets!")
		panic(err)
	}

	rt, err := api.GetRecentTracks("justinbird99")
	if err != nil {
		fmt.Println("Failed to get recent Tracks!")
		panic(err)
	} else if rt.Status != "ok" {
		panic(rt.Error.ErrorMsg)
	}

	token, err := api.GetToken()
	if err != nil {
		fmt.Println("Failed to get token!")
		panic(err)
	} else if token.Status != "ok" {
		panic(token.Error.ErrorMsg)
	}

	fmt.Printf("Token status: %s\n", token.Status)

	fmt.Println("Authorize MostFM to access your LastFM account by clicking this link:")
	fmt.Printf("http://www.last.fm/api/auth/?api_key=%s&token=%s\n", api.APIKey, token.Token)
	fmt.Println("Press enter to continue")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	session, err := api.GetSession(token.Token)
	if err != nil {
		fmt.Println("Failed to get session!")
		panic(err)
	} else if session.Status != "ok" {
		panic(session.Error.ErrorMsg)
	}

	fmt.Println(rt.RecentTracks.User)
	fmt.Println(rt.RecentTracks.Tracks[0].MBID)

	for _, image := range rt.RecentTracks.Tracks[0].Images {
		fmt.Println(image.Size)
		fmt.Println(image.URL)
	}
}
