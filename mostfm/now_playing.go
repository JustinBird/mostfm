package mostfm

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
	"github.com/mattermost/mattermost-server/v6/model"

	"mostfm/lastfm"
)

var NowPlayingForm = apps.Form{
	Title: "Now Playing",
	Icon:  "http://45.76.25.54:4000/static/mostfm.png",
	Fields: []apps.Field{
		{
			Type: "text",
			Name: "Username",
		},
	},
	Submit: apps.NewCall("/now-playing").WithExpand(apps.Expand{
		ActingUser:            apps.ExpandSummary,
		Channel:               apps.ExpandSummary,
	}),
}

//func NowPlayingPost(rt lastfm.RecentTracks) *model.Post {
func NowPlayingPost(channelID string, rt lastfm.RecentTracks) *model.Post {
	post := model.Post {
		ChannelId: channelID,
		Message: "This is a test2!",
	}

	track := rt.Tracks[0]
	authorName := fmt.Sprintf("Now Playing - %s", rt.User)
	if !track.NowPlaying {
		authorName = fmt.Sprintf("Last Played for %s (%s)", rt.User, track.Date.Date)
	}

	attachments := []*model.SlackAttachment {
		{
			AuthorName: authorName,
			AuthorLink: fmt.Sprintf("https://last.fm/user/%s", rt.User),
			Title: "Test",
			Text: "This is a test!!!",
			ImageURL: "https://ia601604.us.archive.org/28/items/mbid-76df3287-6cda-33eb-8e9a-044b5e15ffdd/mbid-76df3287-6cda-33eb-8e9a-044b5e15ffdd-829521842_thumb250.jpg",
		},
	}

	model.ParseSlackAttachment(&post, attachments)
	return &post
}

func NowPlaying(w http.ResponseWriter, req *http.Request, secrets lastfm.Secrets) {
	c := apps.CallRequest{}
	json.NewDecoder(req.Body).Decode(&c)

	if c.Context.Channel == nil {
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to get channel ID!")))
		return
	}

	var username string
	v, ok := c.Values["Username"]
	if ok && v != nil {
		username = v.(string)
	} else {
		err := GetUsername(appclient.AsBot(c.Context), c.Context.ActingUser.Id, &username)
		if err != nil {
			log.Print(err)
			httputils.WriteJSON(w,
				apps.NewErrorResponse(errors.New("Failed to get username!")))
			return
		}

		if username == "" {
			httputils.WriteJSON(w,
				apps.NewErrorResponse(errors.New("No username specified and could not find a registered username! Please specify a username or register with Most.fm.")))
			return
		}
	} 

	rt, err := lastfm.GetRecentTracks(secrets, username)
	if err != nil {
		log.Print(err)
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to get recent tracks!")))
		return
	} else if rt.Status != "ok" {
		log.Print(err)
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New(fmt.Sprintf("Failed to get recent tracks! Error %d: %s.", rt.Error.ErrorCode, rt.Error.ErrorMsg))))
		return
	}

	//post := NowPlayingPost(rt.RecentTracks)
	fmt.Println(c.Context.Channel.Id)
	post := NowPlayingPost(c.Context.Channel.Id, rt.RecentTracks)
	_, err = appclient.AsBot(c.Context).CreatePost(post)
	if err != nil {
		httputils.WriteJSON(w,
			apps.NewErrorResponse(err))
		return
	}

	channelName := c.Context.Channel.DisplayName
	message := fmt.Sprintf("Created a post in %s.", channelName)
	if channelName == "" {
		channelName = c.Context.Channel.Name
		message = "Created a post."
	}
	httputils.WriteJSON(w,
		apps.NewTextResponse(message))
}	
