package mostfm

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
	"github.com/mattermost/mattermost-server/v6/model"
)

func LifeCyclePost(w http.ResponseWriter, req *http.Request, message string) {
	c := apps.CallRequest{}
	json.NewDecoder(req.Body).Decode(&c)

	var post model.Post
	attachments := []*model.SlackAttachment{
		{
			AuthorName: "Most.fm",
			AuthorLink: "https://mostfm.xyz",
			Text:       message,
		},
	}
	model.ParseSlackAttachment(&post, attachments)

	_, err := appclient.AsBot(c.Context).DMPost(c.Context.ActingUser.Id, &post)
	if err != nil {
		httputils.WriteJSON(w, apps.NewErrorResponse(err))
		return
	}

	httputils.WriteJSON(w, apps.NewTextResponse("Created a post in your DM channel."))
}

func (api MostFMAPI) InstallPost(w http.ResponseWriter, req *http.Request) {
	welcome_message := "**Thanks for installing Most.fm!**\n\n" +
		"To begin using Most.fm, make sure that the bot account is invited to the desired team and channel. " +
		"Otherwise the bot will not have the appropriate permissions to post messages. " +
		"To learn more about what Most.fm can do, run the `/mostfm help` command or visit the [Most.fm website](http://mostfm.xyz)."
	LifeCyclePost(w, req, welcome_message)
}

func (api MostFMAPI) UninstallPost(w http.ResponseWriter, req *http.Request) {
	goodbye_message := "**Thanks for trying Most.fm!**\n\n" +
		"We're sorry to hear that you've uninstalled Most.fm. " +
		"If you have had any issues or bugs that you would like to report, please reach out to admin@mostfm.xyz or file a bug on the [GitHub repository](https://github.com/JustinBird/mostfm). "
	LifeCyclePost(w, req, goodbye_message)
}

func (api MostFMAPI) UpdatePost(w http.ResponseWriter, req *http.Request) {
	update_message := "**Most.fm has been updated to a new version!**\n\n" +
		"Visit the [Most.fm website](http://mostfm.xyz) for more about any updates."
	LifeCyclePost(w, req, update_message)
}
