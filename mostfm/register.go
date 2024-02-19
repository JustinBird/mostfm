package mostfm

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
	//"github.com/mattermost/mattermost-server/v6/model"

	"mostfm/lastfm"
)

type MostFMAPI struct {
	LastFM lastfm.LastFMAPI
}

func (api MostFMAPI) Register(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	json.NewDecoder(req.Body).Decode(&c)

	token, err := api.LastFM.GetToken()
	if err != nil {
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New(fmt.Sprintf("Failed to get LastFM Token. Please try again. (%s)", err))))
		return
	}

	header := "Welcome to Most.fm!\n" +
			  "\n" +
			  "Registering with Most.fm allows us to access your Last.fm account. This is only required for some actions. Follow the steps below:\n" +
			  fmt.Sprintf("1. Click [here](http://www.last.fm/api/auth/?api_key=%s&token=%s) to be taken to the Last.fm authorization page\n", api.LastFM.APIKey, token.Token) +
			  "1. If necessary, sign in to your Last.fm account\n" +
			  "1. Click 'YES, ALLOW ACCESS' on the Last.fm authorization page\n" +
			  "1. Click 'Submit' on this MatterMost form\n" +
			  "\n" +
			  "If you successfully authorized Most.fm with your account, you should get a message with your Last.fm user name. " +
			  "Otherwise, you will get an error message describing the issue.\n"

	username_key := fmt.Sprintf("most-fm-username-%s", c.Context.ActingUser.Id)
	var username string
	err = appclient.AsBot(c.Context).KVGet("fm", username_key, &username)
	if err != nil {
		fmt.Printf("%s\n", err)
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to get username!")))
		return
	} else if username == "" {
		header += "\n" +
				  "You currently do not have a registered account on this server."
	} else {
		header += "\n" +
				  fmt.Sprintf("You are currently registered on this server as '%s'", username)
	}
	
	RegisterForm := apps.Form{
		Icon:  "http://45.76.25.54:4000/static/mostfm.png",
		Title: "Register your LastFM account",
		Header: header,
		Submit: apps.NewCall("/validate").WithExpand(apps.Expand{
			ActingUserAccessToken: apps.ExpandAll,
			ActingUser:            apps.ExpandID,
		}).WithState(token.Token),
	}
	httputils.WriteJSON(w,
		apps.NewFormResponse(RegisterForm))
}