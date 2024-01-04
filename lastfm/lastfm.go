package lastfm

import (
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"strings"
)

var LastFMURL = "http://ws.audioscrobbler.com/2.0"

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
	field := Field {"api_sig", signature}
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
