package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"

	"mostfm/lastfm"
)

//go:embed icon.png
var IconData []byte

var Manifest = apps.Manifest{
	AppID: "most-fm",
	Version: "v0.0.1",
	DisplayName: "MostFM",
	Icon: "icon.png",
	HomepageURL: "https://github.com/JustinBird/most-fm",
	RequestedPermissions: []apps.Permission{
		apps.PermissionActAsBot,
	},
	RequestedLocations: []apps.Location{
		apps.LocationCommand,
	},
	Deploy: apps.Deploy{
		HTTP: &apps.HTTP{
			RootURL: "http://mattermost-apps-golang-hello-world:4000",
		},
	},
}

// The details for the App UI bindings
var Bindings = []apps.Binding{
	{
		Location: "/command",
		Bindings: []apps.Binding{
			{
				Icon:        "icon.png",
				Label:       "mostfm",
				Description: "MostFM ",
				Hint:        "[register]",
				Bindings: []apps.Binding{
					{
						Label: "register",
						Hint: "[apikey]",
						Form:  &RegisterForm,
					},
				},
			},
		},
	},
}

var RegisterForm = apps.Form{
	Icon:  "icon.png",
	Title: "Register your LastFM account",
	Submit: apps.NewCall("/register").WithExpand(apps.Expand{
		ActingUserAccessToken: apps.ExpandAll,
		ActingUser:            apps.ExpandID,
	}),
}

// main sets up the http server, with paths mapped for the static assets, the
// bindings callback, and the send function.
func main() {
	// Serve static assets: the manifest and the icon.
	http.HandleFunc("/manifest.json",
		httputils.DoHandleJSON(Manifest))
	http.HandleFunc("/static/icon.png",
		httputils.DoHandleData("image/png", IconData))

	// Bindinings callback.
	http.HandleFunc("/bindings",
		httputils.DoHandleJSON(apps.NewDataResponse(Bindings)))

	http.HandleFunc("/register", Register)

	addr := ":4000" // matches manifest.json
	fmt.Println("Listening on", addr)
	fmt.Println("Use '/apps install http http://mattermost-apps-golang-hello-world" + addr + "/manifest.json' to install the app") // matches manifest.json
	log.Fatal(http.ListenAndServe(addr, nil))
}

func Register(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	fmt.Println("Register called")
	json.NewDecoder(req.Body).Decode(&c)

	var s lastfm.Secrets
	lastfm.GetSecrets(&s)

	var t lastfm.LastFMToken
	lastfm.GetToken(s.APIKey, &t)
	
	message := "Authorize this app to access your LastFM account by clicking this link: "
	message += fmt.Sprintf("http://www.last.fm/api/auth/?api_key=%s&token=%s\n", s.APIKey, t.Token)
	appclient.AsBot(c.Context).DM(c.Context.ActingUser.Id, message)

	httputils.WriteJSON(w,
		apps.NewTextResponse("Created a post in your DM channel."))
}
