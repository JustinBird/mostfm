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
func NowPlayingPost(c apps.CallRequest, rt lastfm.RecentTracks) (*model.Post, error) {
	post := model.Post {
		ChannelId: c.Context.Channel.Id,
		Message: "This is a test2!",
	}

	if len(rt.Tracks) < 1 {
		return &post, errors.New("No track data found!")
	}
	track := rt.Tracks[0]
	authorName := fmt.Sprintf("Now Playing - %s", c.Context.ActingUser.Username)
	if !track.NowPlaying {
		authorName = fmt.Sprintf("Last Played for %s (%s)", c.Context.ActingUser.Username, track.Date.Date)
	}

	attachments := []*model.SlackAttachment {
		{
			AuthorName: authorName,
			AuthorLink: fmt.Sprintf("https://last.fm/user/%s", rt.User),
			Title: track.Name,
			TitleLink: track.URL,
			Text: fmt.Sprintf("**%s** • *%s*", track.Artist.Name, track.Album.Name),
			ImageURL: rt.Tracks[0].Images[3].URL,
			Footer: fmt.Sprintf("%d total scrobbles", rt.Total),
		},
	}

	model.ParseSlackAttachment(&post, attachments)
	return &post, nil
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

	post, err := NowPlayingPost(c, rt.RecentTracks)
	if err != nil {
		httputils.WriteJSON(w,
			apps.NewErrorResponse(err))
		return
	}

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
