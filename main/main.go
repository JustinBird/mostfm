package main

import (
	//"bufio"
	"fmt"
	//"os"

	"mostfm/lastfm"
)

func main() {
	secrets, err := lastfm.GetSecrets("secrets.xml")
	if err != nil {
		fmt.Println("Failed to get secrets!")
		panic(err)
	}

	rt, err := lastfm.GetRecentTracks(secrets, "justinbird99")
	if err != nil {
		fmt.Println("Failed to get recent Tracks!")
		panic(err)
	} else if rt.Status != "ok" {
		panic(rt.Error.ErrorMsg)
	}

	fmt.Println(rt.RecentTracks.User)
	fmt.Println(rt.RecentTracks.Tracks[0].MBID)

	for _, image := range rt.RecentTracks.Tracks[0].Images {
		fmt.Println(image.Size)
		fmt.Println(image.URL)
	}
}