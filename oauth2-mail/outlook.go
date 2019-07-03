package main

// APPS:
// apps.dev.microsoft.com

// REMOVE ACCESS:
// https://portal.office.com/account/#

// Developer:
// https://developer.microsoft.com/en-us/graph/docs/concepts/permissions_reference

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type HeaderData struct {
	SPF         string
	FromDomain  string
	FromName    string
	FromAddress string
	FromSpoof   string
	DKIM        bool
}

type attachment struct {
	OdataType        string `json:"@odata.type"`
	Id               string `json:"id"`
	LastModifiedDate string `json:"lastModifiedDateTime"`
	Name             string `json:"name"`
	ContentType      string `json:"contentType"`
	Size             int32  `json:"size"`
	IsInline         bool   `json:"isInline"`
	ContentId        string `json:"contentId"`
	ContentBytes     string `json:"ContentBytes"`
	ContentLocation  string `json:"ContentLocation"`
}

type attachments struct {
	Value        []attachment `json:"value"`
	OdataContent string       `json:"@odata.content"`
	Raw          []byte       `json:"-"`
}

type emailaddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type from struct {
	EmailAddress emailaddress `json:"emailAddress"`
}

type body struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

type mail struct {
	From             from     `json:"from"`
	Sender           from     `json:"sender"`
	Id               string   `json:"id"`
	Subject          string   `json:"subject"`
	Etag             string   `json:"@odata.etag"`
	Body             body     `json:"body"`
	IsRead           bool     `json:"isRead"`
	Categories       []string `json:"categories"`
	ConversationId   string   `json:"conversationId"`
	ReceivedDateTime string   `json:"receivedDateTime"`
}

type multimail struct {
	Value        []mail `json:"value"`
	DataContent  string `json:"@odata.context"`
	DataNextLink string `json:"@odata.nextLink"`
	Raw          []byte `json:"-"`
}

type AppConfig struct {
	ClientID     string
	ClientSecret string
	RedirectUrl  string
	AuthUrl      string
	TokenUrl     string
	Scope        []string
}

type user struct {
	email   string
	folders []string
}

// Client used for further mail gathering, e.g. Attachments
func mailFilter(client *http.Client, curFolder string, mailBody string) error {
	var err error
	//mailBytes, _ := simplejson.NewJson([]byte(mailBody))
	mailContent := multimail{}
	err = json.Unmarshal([]byte(mailBody), &mailContent)
	if err != nil {
		log.Printf("Can't marshal object: %s\n", err)
		return err
	}

	mailContent.Raw = []byte(mailBody)

	// Loops over all mails in the current folder
	for _, element := range mailContent.Value {
		log.Println(element.Subject)
	}

	return nil
}

// Sets the office365 config based on a specified token
func setConfig(config AppConfig) (*http.Client, error) {
	ctx := context.Background()

	// RedirectURL needs to be in manifest["replyUrls"]
	conf := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scope,
		RedirectURL:  config.RedirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthUrl,
			TokenURL: config.TokenUrl,
		},
	}

	url := conf.AuthCodeURL("state", oauth2.SetAuthURLParam("resource", "https://graph.microsoft.com"), oauth2.SetAuthURLParam("access_type", "offline"))

	var code string

	log.Printf("Visit the URL for the auth dialog: \n%v\n\n", url)

	// Handles the server callback, listening on port 8000
	codechannel := make(chan string)
	go func() {
		port := ":8000"

		http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
			tmpcode := request.URL.Query().Get("code")
			if len(tmpcode) < 100 {
				return
			} else {
				codechannel <- tmpcode
			}
		})

		// FIX - might cause errors not being printed
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Printf("%s\n", err)
		}
	}()

	for {
		code = <-codechannel
		access_token, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Printf("Wrong code: %s", code)
			continue
		}

		client := conf.Client(ctx, access_token)
		close(codechannel)
		return client, nil
	}
}

func makeRequests(client *http.Client, mailAmount string, curUser user) {
	var requestUrl string

	//requestUrl := "https://graph.microsoft.com/v1.0/me/mailfolders"
	requestUrl = fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/mailfolders", curUser.email)

	// Loops over mail for specific folder based on childIds
	for _, newId := range curUser.folders {
		endpointUrl := fmt.Sprintf("%s/%s/messages?$select=categories,subject,body,from,attachments,isRead,receivedDateTime,conversationId&$top=%s&$orderby=receivedDateTime%sDESC", requestUrl, newId, mailAmount, "%20")

		endpointRet, err := client.Get(endpointUrl)

		if err != nil {
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(endpointRet.Body)
		newStr := buf.String()

		// Goes through all mail
		err = mailFilter(client, newId, newStr)
		if err != nil {
			continue
		}
	}
}

func findMailboxes(folders []string) []user {
	users := []user{}
	for _, line := range folders {
		if len(line) == 0 {
			continue
		}

		if !strings.Contains(line, ":") {
			continue
		}

		linesplit := strings.Split(line, ":")
		if len(linesplit) > 2 {
			log.Printf("Error in line: %s", line)
			continue
		}

		if !strings.Contains(linesplit[0], "@") {
			log.Printf("Error in line: %s", line)
			continue
		}

		curuser := user{
			email: linesplit[0],
		}

		if strings.Contains(linesplit[1], "/") {
			for _, item := range strings.Split(linesplit[1], "/") {
				curuser.folders = append(curuser.folders, item)
			}
		} else {
			curuser.folders = append(curuser.folders, linesplit[1])
		}

		users = append(users, curuser)
	}

	return users
}

func checkConfig(appconfig AppConfig) error {
	if len(appconfig.ClientID) == 0 {
		return errors.New("ClientID configuration can't be empty.")
	}
	if len(appconfig.ClientSecret) == 0 {
		return errors.New("ClientSecret configuration can't be empty.")
	}
	if len(appconfig.RedirectUrl) == 0 {
		return errors.New("RedirectUrl configuration can't be empty.")
	}
	if len(appconfig.AuthUrl) == 0 {
		return errors.New("AuthUrl configuration can't be empty.")
	}
	if len(appconfig.TokenUrl) == 0 {
		return errors.New("TokenUrl configuration can't be empty.")
	}
	if len(appconfig.Scope) == 0 {
		return errors.New("Scope configuration can't be empty.")
	}

	return nil
}

// Configure your app
func main() {
	// Interval between new requests to same folder
	secondInterval := 120
	// Amount of mail to look for in each folder every <interval> seconds
	// Max is 1000
	mailAmount := "1000"

	// Configuration of oauth2 app
	appconfig := AppConfig{
		ClientID:     "",
		ClientSecret: "",
		RedirectUrl:  "",
		AuthUrl:      "",
		TokenUrl:     "",
		Scope:        []string{},
	}

	// Configure folders to read from
	folders := []string{
		"example@example.com:inbox",
		"example@example.com:inbox/subfolder/subsubfolder",
	}

	err := checkConfig(appconfig)
	if err != nil {
		panic(err)
	}

	// Parses the folders defined above if valid
	users := findMailboxes(folders)
	if len(users) == 0 {
		panic("No folders found")
	}

	// Builds the oauth query
	client, err := setConfig(appconfig)
	if err != nil {
		panic(err)
	}

	log.Printf("Got callback! Looking for %d mail on a one minute period in each of the %d mailboxes.", mailAmount, len(users))
	for {
		for _, mailfolder := range users {
			makeRequests(client, mailAmount, mailfolder)
		}

		time.Sleep(time.Duration(secondInterval) * time.Second)
	}
}
