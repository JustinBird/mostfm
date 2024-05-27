package mostfm

import (
	_ "embed"
	"fmt"

	"github.com/mattermost/mattermost-plugin-apps/apps/appclient"
)

var keyPrefix string = "fm"

func usernameKeyName(id string) string {
	return fmt.Sprintf("most-fm-username-%s", id)
}

func sessionKeyName(id string) string {
	return fmt.Sprintf("most-fm-session-%s", id)
}

func GetUsername(c *appclient.Client, id string) (username string, err error) {
	usernameKey := usernameKeyName(id)
	err = c.KVGet(keyPrefix, usernameKey, &username)
	return
}

func GetSession(c *appclient.Client, id string) (session string, err error) {
	sessionKey := sessionKeyName(id)
	err = c.KVGet(keyPrefix, sessionKey, &session)
	return
}

func SetUsername(c *appclient.Client, id string, username string) error {
	usernameKey := usernameKeyName(id)
	_, err := c.KVSet(keyPrefix, usernameKey, username)
	return err
}

func SetSession(c *appclient.Client, id string, session string) error {
	sessionKey := sessionKeyName(id)
	_, err := c.KVSet(keyPrefix, sessionKey, session)
	return err
}
