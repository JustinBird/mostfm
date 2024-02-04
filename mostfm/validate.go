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

func Validate(w http.ResponseWriter, req *http.Request, secrets lastfm.Secrets) {
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