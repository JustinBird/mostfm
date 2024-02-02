package main

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
				Hint:        "[register]",
				Bindings: []apps.Binding{
					{
						Label: "register",
						Submit: apps.NewCall("/register").WithExpand(apps.Expand{
							ActingUserAccessToken: apps.ExpandAll,
							ActingUser:            apps.ExpandID,
						}),
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

	http.HandleFunc("/register", Register)
	http.HandleFunc("/validate", Validate)

	addr := ":4000" // matches manifest.json
	fmt.Println("Listening on", addr)
	fmt.Println("Use '/apps install http http://mattermost-apps-golang-hello-world" + addr + "/manifest.json' to install the app") // matches manifest.json
	log.Fatal(http.ListenAndServe(addr, nil))
}

func Register(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	fmt.Println("Register called")
	json.NewDecoder(req.Body).Decode(&c)

	token, err := lastfm.GetToken(secrets.APIKey)
	if err != nil {
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New(fmt.Sprintf("Failed to get LastFM Token. Please try again. (%s)", err))))
		return
	}

	header := "Welcome to Most.fm!\n" +
			  "\n" +
			  "Registering with Most.fm allows us to access your Last.fm account. This is only required for some actions. Follow the steps below:\n" +
			  fmt.Sprintf("1. Click [here](http://www.last.fm/api/auth/?api_key=%s&token=%s) to be taken to the Last.fm authorization page\n", secrets.APIKey, token.Token) +
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

func Validate(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	fmt.Println("Validate called")
	json.NewDecoder(req.Body).Decode(&c)

	token, ok := c.Call.State.(string)
	if !ok || len(token) == 0 {
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Token was not included in response! Please try again.")))
		return
	}

	session, err := lastfm.GetSession(secrets, token)
	if err != nil {
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to get session token! Please try again.")))
		return
	} else if session.Status != "ok" {
		fmt.Println("%s, %d\n", session.Error.ErrorMsg, session.Error.ErrorCode)
		if session.Error.ErrorCode == 14 {
			httputils.WriteJSON(w,
				apps.NewErrorResponse(errors.New("Your token has not been authorized. Please click the link above and allow Most.fm to access your account.")))
		} else {
			httputils.WriteJSON(w,
				apps.NewErrorResponse(errors.New(fmt.Sprintf("Failed to get session token! Error %d: %s.", session.Error.ErrorCode, session.Error.ErrorMsg))))
		}
		return
	}

	fmt.Printf("Acting User ID: %s\n", c.Context.ActingUser.Id)
	fmt.Printf("Acting User username: %s\n", c.Context.ActingUser.Nickname)
	username_key := fmt.Sprintf("most-fm-username-%s", c.Context.ActingUser.Id)
	_, err = appclient.AsBot(c.Context).KVSet("fm", username_key, session.Name)
	if err != nil {
		fmt.Printf("%s\n", err)
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to store username!")))
		return
	}
	session_key := fmt.Sprintf("most-fm-session-%s", c.Context.ActingUser.Id)
	_, err = appclient.AsBot(c.Context).KVSet("fm", session_key, session.Key)
	if err != nil {
		fmt.Printf("%s\n", err)
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to store session key!")))
		return
	}

	httputils.WriteJSON(w,
		apps.NewTextResponse(fmt.Sprintf("Successfully registered with Most.fm under '%s'", session.Name)))
}
