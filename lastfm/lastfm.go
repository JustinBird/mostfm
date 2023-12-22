package lastfm

import (
	"fmt"
	"net/http"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io"
	"os"
	"strings"
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

func GetSecrets(secrets_path string) (Secrets, error) {
	var s Secrets
	data, err := os.ReadFile(secrets_path)
	if err != nil {
		return s, err
	}

	xml.Unmarshal(data, &s)
	return s, nil
}

func GetToken(apikey string) (LastFMToken, error) {
	var t LastFMToken
	fields := []Field {
		{"api_key", apikey},
		{"method",  "auth.getToken"},
	}
	resp, err := http.Get(createURL(fields))
	if err != nil {
		fmt.Println("Failed to get token!")
		return t, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return t, err
	}
	fmt.Printf(string(body))
	err = xml.Unmarshal(body, &t)
	if err != nil {
		return t, err
	}

	if  t.Status != "ok" {
		fmt.Printf("Bad status when getting token: %s\n", t.Status)
	}
	return t, nil
}

func GetSession(secrets Secrets, token string) (LastFMSession, error) {
	var s LastFMSession
	fields := []Field {
		{"api_key", secrets.APIKey},
		{"method", "auth.getSession"},
		{"token", token},
	}
	createSignature(&fields, secrets.Secret)

	resp, err := http.Get(createURL(fields))
	if err != nil {
		fmt.Println("Failed to get session!")
		return s, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	fmt.Printf(string(body))
	xml.Unmarshal(body, &s)

	if  s.Status != "ok" {
		fmt.Printf("Bad status when getting token: %s\n", s.Status)
	}
	fmt.Println(s.Name)
	fmt.Println(s.Key)
	fmt.Println(s.Key)
	return s, nil
}
