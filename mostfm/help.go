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
	"github.com/mattermost/mattermost-server/v6/model"
)

var HelpForm = apps.Form{
	Title: "Help",
	Icon:  "http://mostfm.xyz:4000/static/mostfm.png",
	Submit: apps.NewCall("/help").WithExpand(apps.Expand{
		ActingUser: apps.ExpandSummary,
		Channel:    apps.ExpandSummary,
	}),
}

func (api MostFMAPI) Help(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	json.NewDecoder(req.Body).Decode(&c)

	if c.Context.Channel == nil {
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to get channel ID!")))
		return
	}

	post := model.Post{
		ChannelId: c.Context.Channel.Id,
	}

	helpMessage := "[Most.fm](www.mostfm.xyz) is a bot for MatterMost that allows you to share and interact with Last.fm. " +
		"A sub-command can be run by typing `/mostfm <sub-command>` in the message box. " +
		"Some sub-commands take parameters which can be given by typing `/mostfm <sub-command> --<parameter> <value>`. " +
		"The following sub-commands are available:" +
		"\n\n" +
		"| Sub-command | Action | Parameters |\n" +
		"| :---------- | :----- | :--------- |\n" +
		"| `help`        | Posts this menu | None    |\n" +
		"| `register`    | Registers you with Most.fm | None |\n" +
		"| `now-playing` | Posts your lastest or currently playing track | - `username`: The username to get the latest track from. Required unless you have registered with Most.fm. |\n" +
		"\n\n" +
		"For more info visit www.mostfm.xyz"

	attachments := []*model.SlackAttachment{
		{
			AuthorName: "Most.fm",
			AuthorLink: "https://mostfm.xyz",
			Text:       helpMessage,
		},
	}

	model.ParseSlackAttachment(&post, attachments)

	_, err := appclient.AsBot(c.Context).CreatePost(&post)
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
