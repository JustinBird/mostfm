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
	AppID:       "most-fm",
	Version:     "v0.0.1",
	DisplayName: "Most.fm",
	Icon:        "mostfm.png",
	HomepageURL: "https://mostfm.xyz",
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
	OnInstall: &apps.Call{
		Path: "/install",
		Expand: &apps.Expand{
			ActingUser: apps.ExpandID,
		},
	},
	OnUninstall: &apps.Call{
		Path: "/uninstall",
		Expand: &apps.Expand{
			ActingUser: apps.ExpandID,
		},
	},
	OnVersionChanged: &apps.Call{
		Path: "/version_changed",
		Expand: &apps.Expand{
			ActingUser: apps.ExpandID,
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
						Form:  &mostfm.NowPlayingForm,
					},
				},
			},
		},
	},
}

var api mostfm.MostFMAPI

// main sets up the http server, with paths mapped for the static assets, the
// bindings callback, and the send function.
func main() {
	var err error
	api.LastFM, err = lastfm.NewAPIFromFile("secrets.xml")
	if err != nil {
		log.Fatal("Failed to get secrets!")
	}
	fmt.Printf("Using API key: %s\n", api.LastFM.APIKey)

	// Serve static assets: the manifest and the icon.
	http.HandleFunc("/manifest.json",
		httputils.DoHandleJSON(Manifest))
	http.HandleFunc("/static/mostfm.png",
		httputils.DoHandleData("image/png", IconData))

	// Bindings callback.
	http.HandleFunc("/bindings",
		httputils.DoHandleJSON(apps.NewDataResponse(Bindings)))

	http.HandleFunc("/install", api.InstallPost)
	http.HandleFunc("/uninstall", api.UninstallPost)
	http.HandleFunc("/version_changed", api.UpdatePost)
	http.HandleFunc("/register", api.Register)
	http.HandleFunc("/validate", api.Validate)
	http.HandleFunc("/now-playing", api.NowPlaying)

	addr := ":4000" // matches manifest.json
	fmt.Println("Listening on", addr)
	fmt.Println("Use '/apps install http http://mattermost-apps-golang-hello-world" + addr + "/manifest.json' to install the app") // matches manifest.json
	log.Fatal(http.ListenAndServe(addr, nil))
}
