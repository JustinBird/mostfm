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
	//"github.com/mattermost/mattermost-server/v6/model"

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
		ActingUserAccessToken: apps.ExpandAll,
		ActingUser:            apps.ExpandID,
	}),
}

func NowPlaying(w http.ResponseWriter, req *http.Request, secrets lastfm.Secrets) {
	c := apps.CallRequest{}
	fmt.Println("Now playing called")
	json.NewDecoder(req.Body).Decode(&c)

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
	} else if rt.Status != "ok" {
		log.Print(err)
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New(fmt.Sprintf("Failed to get recent tracks! Error %d: %s.", rt.Error.ErrorCode, rt.Error.ErrorMsg))))
	}

	fmt.Println("End of now playing %s", username)
	httputils.WriteJSON(w,
		apps.NewTextResponse(fmt.Sprintf("%s %s\n", username,  rt.RecentTracks.Tracks[0].Name)))
}	
