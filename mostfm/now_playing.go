package mostfm

import (
	_ "embed"
	"encoding/json"
	//"errors"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	//"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
	//"github.com/mattermost/mattermost-server/v6/model"

	//"mostfm/lastfm"
)

func NowPlaying(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	fmt.Println("Now playing called")
	json.NewDecoder(req.Body).Decode(&c)

	httputils.WriteJSON(w,
		apps.NewTextResponse("Now playing"))
}	
