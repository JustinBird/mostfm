package lastfm

import(
	"fmt"
	"encoding/xml"
)

type Secrets struct {
	APIKey string `xml:"apikey"`
	Secret string `xml:"secret"`
}

type Field struct {
	Key   string
	Value string
}

type LastFMError struct {
	XMLName   xml.Name `xml:"error"`
	ErrorMsg  string   `xml:",chardata"`
	ErrorCode int      `xml:"code,attr"`
}

func (e LastFMError) String() string {
	return fmt.Sprintf("%s (%d)", e.ErrorMsg, e.ErrorCode)
}

type LastFMToken struct {
	XMLName xml.Name    `xml:"lfm"`
	Status  string      `xml:"status,attr"`
	Token   string      `xml:"token"`
	Error   LastFMError `xml:"error"`
}

func (t LastFMToken) String() string {
	if (t.Status == "ok") {
		return fmt.Sprintf("Token %s (%s)", t.Status, t.Token)
	} else {
		return t.Error.String()
	}
}

type LastFMSession struct {
	XMLName    xml.Name    `xml:"lfm"`
	Status     string      `xml:"status,attr"`
	Name       string      `xml:"session>name"`
	Key        string      `xml:"session>key"`
	Subscriber int         `xml:"session>subscriber"`
	Error      LastFMError `xml:"error"`
}

type Date struct {
	XMLName   xml.Name `xml:"date"`
	Date      string   `xml:",chardata"`
	TimeStamp int      `xml:"uts,attr"`
}

func (d Date) String() string {
	return d.Date
}

type Artist struct {
	XMLName xml.Name `xml:"artist"`
	Name  string   `xml:",chardata"`
	MBID    string   `xml:"mbid,attr"`
}

func (a Artist) String() string {
	return a.Name
}

type Album struct {
	XMLName xml.Name `xml:"album"`
	Name    string   `xml:",chardata"`
	MBID    string   `xml:"mbid,attr"`
}

type Image struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:",chardata"`
	Size    string   `xml:"size,attr"`
}
type Track struct {
	XMLName    xml.Name  `xml:"track"`
	NowPlaying bool      `xml:"nowplaying,attr"`
	Artist     Artist    `xml:"artist"`
	Name       string    `xml:"name"`
	MBID       string    `xml:"mbid"`
	Album      Album     `xml:"album"`
	URL        string    `xml:"url"`
	Streamable int       `xml:"streamable"`
	Date       Date      `xml:"date"`
	Images     []Image   `xml:"image"`
}

func (t Track) String() string {
	if t.NowPlaying {
		return fmt.Sprintf("Currently playing: %s - \"%s\"", t.Artist, t.Name)
	} else {
		return fmt.Sprintf("Last played: %s - \"%s\" on %s\n", t.Artist, t.Name, t.Date)
	}
}

type RecentTracks struct {
	XMLName    xml.Name `xml:"recenttracks"`
	User       string   `xml:"user,attr"`
	Page       int      `xml:"page,attr"`
	PerPage    int      `xml:"perPage,attr"`
	TotalPages int      `xml:"totalPages,attr"`
	Total      int      `xml:"total,attr"`
	Tracks     []Track  `xml:"track"`
}

type LastFMRecentTracks struct {
	XMLName        xml.Name    `xml:"lfm"`
	Status        string       `xml:"status,attr"`
	RecentTracks  RecentTracks `xml:"recenttracks"`
	Error         LastFMError  `xml:"error"`
}

type LastFMResponse interface {
	*LastFMToken | *LastFMSession | *LastFMRecentTracks
}