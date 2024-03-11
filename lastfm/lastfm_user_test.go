package lastfm

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetRecentTracks(t *testing.T) {
	var api = LastFMAPI{APIKey: "0123456789", Secret: "abcdefg"}
	tests := []struct {
		API                  LastFMAPI
		Output               string
		HTTPCode             int
		ExpectedRecentTracks LastFMRecentTracks
		ExpectedError        error
	}{
		{ // Normal operation
			API: api,
			Output: `
			<?xml version="1.0" encoding="UTF-8"?>
			<lfm status="ok">
			  <recenttracks user="mostfm" page="1" perPage="10" totalPages="10" total="100">
				<track>
				  <artist mbid="artistmbid">artist</artist>
				  <name>name</name>
				  <streamable>0</streamable>
				  <mbid>mbid</mbid>
				  <album mbid="albummbid">album</album>
				  <url>url</url>
				  <image size="small">small</image>
				  <image size="medium">medium</image>
				  <image size="large">large</image>
				  <image size="extralarge">extralarge</image>
				  <date uts="uts">date</date>
				</track>
			  </recenttracks>
			</lfm>`,
			HTTPCode: http.StatusOK,
			ExpectedRecentTracks: LastFMRecentTracks{
				XMLName: xml.Name{Local: "lfm"},
				Status:  "ok",
				RecentTracks: RecentTracks{
					XMLName:    xml.Name{Local: "recenttracks"},
					User:       "mostfm",
					Page:       1,
					PerPage:    10,
					TotalPages: 10,
					Total:      100,
					Tracks: []Track{
						{
							XMLName:    xml.Name{Local: "track"},
							NowPlaying: false,
							Artist: Artist{
								XMLName: xml.Name{Local: "artist"},
								Name:    "artist",
								MBID:    "artistmbid",
							},
							Name: "name",
							MBID: "mbid",
							Album: Album{
								XMLName: xml.Name{Local: "album"},
								Name:    "album",
								MBID:    "albummbid",
							},
							URL:        "url",
							Streamable: 0,
							Date: Date{
								XMLName:   xml.Name{Local: "date"},
								Date:      "date",
								TimeStamp: "uts",
							},
							Images: []Image{
								{XMLName: xml.Name{Local: "image"}, URL: "small", Size: "small"},
								{XMLName: xml.Name{Local: "image"}, URL: "medium", Size: "medium"},
								{XMLName: xml.Name{Local: "image"}, URL: "large", Size: "large"},
								{XMLName: xml.Name{Local: "image"}, URL: "extralarge", Size: "extralarge"},
							},
						},
					},
				},
			},
			ExpectedError: nil,
		},
	}

	for i, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.HTTPCode)
			w.Write([]byte(test.Output))
		}))

		LastFMURL = server.URL
		rt, err := test.API.GetRecentTracks("mostfm")
		if !errors.Is(err, test.ExpectedError) {
			t.Errorf("Test #%d: Expected error message did not match!\nExpected: %s\nActual: %s", i, test.ExpectedError, err)
		}

		if !reflect.DeepEqual(rt, test.ExpectedRecentTracks) {
			t.Errorf("Test #%d: LastFM recent tracks did not match expected!", i)
		}
	}
}
