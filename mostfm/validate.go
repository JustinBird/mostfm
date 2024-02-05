package mostfm

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"log"

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
		if session.Error.ErrorCode == 14 {
			httputils.WriteJSON(w,
				apps.NewErrorResponse(errors.New("Your token has not been authorized. Please click the link above and allow Most.fm to access your account.")))
		} else {
			httputils.WriteJSON(w,
				apps.NewErrorResponse(errors.New(fmt.Sprintf("Failed to get session token! Error %d: %s.", session.Error.ErrorCode, session.Error.ErrorMsg))))
		}
		return
	}

	err = SetUsername(appclient.AsBot(c.Context), c.Context.ActingUser.Id, session.Name)
	if err != nil {
		log.Print(err)
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to store username!")))
		return
	}

	err = SetSession(appclient.AsBot(c.Context), c.Context.ActingUser.Id, session.Key)
	if err != nil {
		log.Print(err)
		httputils.WriteJSON(w,
			apps.NewErrorResponse(errors.New("Failed to store session key!")))
		return
	}

	httputils.WriteJSON(w,
		apps.NewTextResponse(fmt.Sprintf("Successfully registered with Most.fm under '%s'", session.Name)))
}