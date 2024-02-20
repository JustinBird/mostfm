package lastfm

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var LastFMURL = "http://ws.audioscrobbler.com/2.0"

var ErrFileRead = errors.New("failed to read file")
var ErrHTTPCall = errors.New("failed to make Last.fm HTTP call")
var ErrHTTPCode = errors.New("bad HTTP Code")
var ErrReadBody = errors.New("failed to read body")
var ErrXMLParse = errors.New("bad XML data")
var ErrLastFMStatus = errors.New("bad Last.fm status")

func NewAPI(api_key, secret string) LastFMAPI {
	return LastFMAPI{
		APIKey: api_key,
		Secret: secret,
	}
}

func NewAPIFromFile(path string) (LastFMAPI, error) {
	var api LastFMAPI
	data, err := os.ReadFile(path)
	if err != nil {
		return api, errors.Join(ErrFileRead, err)
	}

	err = xml.Unmarshal(data, &api)
	if err != nil {
		return api, errors.Join(ErrXMLParse, err)
	}

	return api, nil
}

func createSignature(fields *[]Field, shared_secret string) {
	var data strings.Builder
	for _, field := range *fields {
		data.WriteString(field.Key)
		data.WriteString(field.Value)
	}
	data.WriteString(shared_secret)

	bytes := []byte(data.String())
	hash := md5.Sum(bytes)
	signature := hex.EncodeToString(hash[:])
	field := Field{"api_sig", signature}
	*fields = append(*fields, field)
}

func createURL(fields []Field) string {
	var url strings.Builder
	url.WriteString(LastFMURL + "/?")
	for _, field := range fields {
		url.WriteString(fmt.Sprintf("%s=%s&", field.Key, field.Value))
	}
	return url.String()
}

func LastFMCall[T LastFMResponse](fields *[]Field, r T) error {
	resp, err := http.Get(createURL(*fields))
	if err != nil {
		return errors.Join(err, ErrHTTPCall)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Join(err, ErrReadBody)
	}

	// If there's no body to parse then fail, otherwise try to parse it
	if resp.StatusCode != 200 && len(body) == 0 {
		return fmt.Errorf("%w HTTP status code %d", ErrHTTPCode, resp.StatusCode)
	}

	err = xml.Unmarshal(body, r)
	if err != nil {
		return errors.Join(err, ErrXMLParse)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%w HTTP status code %d", ErrHTTPCode, resp.StatusCode)
	}

	return nil
}
