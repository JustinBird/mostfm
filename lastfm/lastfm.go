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

type LastFMToken struct {
	XMLName xml.Name `xml:"lfm"`
	Status  string   `xml:"status,attr"`
	Token   string   `xml:"token"` 
}

type LastFMSession struct {
	XMLName    xml.Name `xml:"lfm"`
	Status     string   `xml:"status,attr"`
	Name       string   `xml:"session>name"`
	Key        string   `xml:"session>key"`
	Subscriber int      `xml:"session>subscriber"`
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
	fmt.Printf("My signature %s\n", signature)
	field := Field {"api_sig", signature}
	*fields = append(*fields, field)
}

func createURL(fields []Field) string {
	var url strings.Builder
	url.WriteString("http://ws.audioscrobbler.com/2.0/?")
	for _, field := range fields {
		fmt.Println(field.Key)
		fmt.Println(field.Value)
		url.WriteString(fmt.Sprintf("%s=%s&", field.Key, field.Value))
	}
	fmt.Println(url.String())
	return url.String()
}

func GetSecrets(s *Secrets) {
	data, err := os.ReadFile("secrets.xml")
	if err != nil {
		fmt.Println("Failed to read secrets file!")
		panic(err)
	}

	xml.Unmarshal(data, s)
}

func GetToken(apikey string, t *LastFMToken) {
	fields := []Field {
		{"api_key", apikey},
		{"method",  "auth.getToken"},
	}
	resp, err := http.Get(createURL(fields))
	if err != nil {
		fmt.Println("Failed to get token!")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	fmt.Printf(string(body))
	xml.Unmarshal(body, t)

	if  t.Status != "ok" {
		fmt.Printf("Bad status when getting token: %s\n", t.Status)
	}
}

func GetSession(secrets Secrets, token string, s *LastFMSession) {
	fields := []Field {
		{"api_key", secrets.APIKey},
		{"method", "auth.getSession"},
		{"token", token},
	}
	createSignature(&fields, secrets.Secret)

	resp, err := http.Get(createURL(fields))
	if err != nil {
		fmt.Println("Failed to get session!")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("HTTP status code %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	fmt.Printf(string(body))
	xml.Unmarshal(body, s)

	if  s.Status != "ok" {
		fmt.Printf("Bad status when getting token: %s\n", s.Status)
	}
	fmt.Println(s.Name)
	fmt.Println(s.Key)
	fmt.Println(s.Key)
}


