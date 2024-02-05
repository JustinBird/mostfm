package main

import (
	"bufio"
	"fmt"
	"os"

	"mostfm/lastfm"
)

func main() {
	secrets, err := lastfm.GetSecrets("secrets.xml")
	if err != nil {
		fmt.Println("Failed to get secrets!")
		panic(err)
	}

	token, err := lastfm.GetToken(secrets.APIKey)
	if err != nil {
		fmt.Println("Failed to get token!")
		panic(err)
	} else if token.Status != "ok" {
		panic(token.Error.ErrorMsg)
	}

	fmt.Printf("Token status: %s\n", token.Status)
	
	fmt.Println("Authorize MostFM to access your LastFM account by clicking this link:")
	fmt.Printf("http://www.last.fm/api/auth/?api_key=%s&token=%s\n", secrets.APIKey, token.Token)
	fmt.Println("Press enter to continue")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	session, err := lastfm.GetSession(secrets, token.Token)
	if err != nil {
		fmt.Println("Failed to get session!")
		panic(err)
	} else if session.Status != "ok" {
		panic(session.Error.ErrorMsg)
	}

	fmt.Printf("Session status: %s\n", session.Status)
	fmt.Printf("Error! %s: %d\n", session.Error.ErrorMsg, session.Error.ErrorCode)

	rt, err := lastfm.GetTracks(secrets, session.Name)
	if err != nil {
		fmt.Println("Failed to get recent Tracks!")
		panic(err)
	} else if rt.Status != "ok" {
		panic(rt.Error.ErrorMsg)
	}

	for _, track := range rt.RecentTracks.Tracks {
		fmt.Printf("%s", track)
	}
}