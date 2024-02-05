package main

import (
	_ "embed"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	//"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
	//"github.com/mattermost/mattermost-server/v6/model"

	"mostfm/lastfm"
	"mostfm/mostfm"
)

//go:embed mostfm.png
var IconData []byte

var Manifest = apps.Manifest{
	AppID: "most-fm",
	Version: "v0.0.1",
	DisplayName: "Most.fm",
	Icon: "mostfm.png",
	HomepageURL: "https://github.com/JustinBird/most-fm",
	RequestedPermissions: []apps.Permission{
		apps.PermissionActAsBot,
	},
	RequestedLocations: []apps.Location{
		apps.LocationCommand,
		apps.LocationInPost,
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
				Icon:        "http://45.76.25.54:4000/static/mostfm.png",
				Label:       "mostfm",
				Description: "Most.fm",
				Hint:        "[register|now-playing]",
				Bindings: []apps.Binding{
					{
						Label: "register",
						Submit: apps.NewCall("/register").WithExpand(apps.Expand{
							ActingUserAccessToken: apps.ExpandAll,
							ActingUser:            apps.ExpandID,
						}),
					},
					{
						Label: "now-playing",
						Form: &mostfm.NowPlayingForm,
					},
				},
			},
		},
	},
}

var secrets = lastfm.Secrets{
	APIKey: "",
	Secret: "",
}

// main sets up the http server, with paths mapped for the static assets, the
// bindings callback, and the send function.
func main() {
	var err error
	secrets, err = lastfm.GetSecrets("secrets.xml")
	if err != nil {
		log.Fatal("Failed to get secrets!")
	}
	fmt.Printf("Using API key: %s\n", secrets.APIKey)

	// Serve static assets: the manifest and the icon.
	http.HandleFunc("/manifest.json",
		httputils.DoHandleJSON(Manifest))
	http.HandleFunc("/static/mostfm.png",
		httputils.DoHandleData("image/png", IconData))

	// Bindinings callback.
	http.HandleFunc("/bindings",
		httputils.DoHandleJSON(apps.NewDataResponse(Bindings)))

	http.HandleFunc("/register",    func(w http.ResponseWriter, r *http.Request) { mostfm.Register(w, r, secrets)   })
	http.HandleFunc("/validate",    func(w http.ResponseWriter, r *http.Request) { mostfm.Validate(w, r, secrets)   })
	http.HandleFunc("/now-playing", func(w http.ResponseWriter, r *http.Request) { mostfm.NowPlaying(w, r, secrets) })

	addr := ":4000" // matches manifest.json
	fmt.Println("Listening on", addr)
	fmt.Println("Use '/apps install http http://mattermost-apps-golang-hello-world" + addr + "/manifest.json' to install the app") // matches manifest.json
	log.Fatal(http.ListenAndServe(addr, nil))
}