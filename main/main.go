package main

import (
	"bufio"
	"fmt"
	"os"
	

	"mostfm/lastfm"
)

func main() {
	s, err := lastfm.GetSecrets("secrets.xml")
	if err != nil {
		fmt.Println("Failed to get secrets!")
		panic(err)
	}

	var t lastfm.LastFMToken
	lastfm.GetToken(s.APIKey, &t)
	
	fmt.Println("Authorize this app to access your LastFM account by clicking this link:")
	fmt.Printf("http://www.last.fm/api/auth/?api_key=%s&token=%s\n", s.APIKey, t.Token)
	fmt.Println("Press enter to continue")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	var session lastfm.LastFMSession
	lastfm.GetSession(s, t.Token, &session)
}