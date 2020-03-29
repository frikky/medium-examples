package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// The message we get from Messenger
type InputMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
				Nlp  struct {
					Entities struct {
						Sentiment []struct {
							Confidence float64 `json:"confidence"`
							Value      string  `json:"value"`
						} `json:"sentiment"`
						Greetings []struct {
							Confidence float64 `json:"confidence"`
							Value      string  `json:"value"`
						} `json:"greetings"`
					} `json:"entities"`
					DetectedLocales []struct {
						Locale     string  `json:"locale"`
						Confidence float64 `json:"confidence"`
					} `json:"detected_locales"`
				} `json:"nlp"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

// The recipient of our message
type Recipient struct {
	ID string `json:"id"`
}

// The message to send it its basic
type Message struct {
	Text string `json:"text,omitempty"`
}

type Button struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Payload string `json:"payload,omitempty"`
	URL     string `json:"url,omitempty"`
}

type Element struct {
	Title         string        `json:"title,omitempty"`
	Subtitle      string        `json:"subtitle,omitempty"`
	ImageURL      string        `json:"image_url,omitempty"`
	DefaultAction DefaultAction `json:"default_action,omitempty"`
	Buttons       []Button      `json:"buttons,omitempty"`
}

type DefaultAction struct {
	Type                string `json:"type,omitempty"`
	URL                 string `json:"url,omitempty"`
	WebViewHeightRation string `json:"webview_height_ratio,omitempty"`
}

// The attachment to send (custom)
type Attachment struct {
	Attachment struct {
		Type    string `json:"type,omitempty"`
		Payload struct {
			TemplateType string    `json:"template_type,omitempty"`
			Elements     []Element `json:"elements,omitempty"`
		} `json:"payload,omitempty"`
	} `json:"attachment,omitempty"`
}

// Full response
type ResponseAttachment struct {
	Recipient Recipient  `json:"recipient"`
	Message   Attachment `json:"message,omitempty"`
}

// Full response
type ResponseMessage struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message,omitempty"`
}

func handleWebhook(resp http.ResponseWriter, request *http.Request) {
	secretKey := "secret_token"
	if request.Method == "GET" {
		u, _ := url.Parse(request.RequestURI)
		values, _ := url.ParseQuery(u.RawQuery)
		token := values.Get("hub.verify_token")
		if token == secretKey {
			resp.WriteHeader(200)
			resp.Write([]byte(values.Get("hub.challenge")))
			return
		}

		resp.WriteHeader(400)
		resp.Write([]byte(`Bad token`))
		return
	}

	// Anything that reaches here is POST.
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Failed parsing body: %s", err)
		resp.WriteHeader(400)
		resp.Write([]byte("An error occurred"))
		return
	}

	// Parse message into the Message struct
	var message InputMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Printf("Failed unmarshalling message: %s", err)
		resp.WriteHeader(400)
		resp.Write([]byte("An error occurred"))
		return
	}

	// Find messages
	//log.Printf("Message: %#v", message)
	log.Printf("%#v", message)
	for _, entry := range message.Entry {
		if len(entry.Messaging) == 0 {
			log.Printf("No messages")
			resp.WriteHeader(400)
			resp.Write([]byte("An error occurred"))
			return
		}

		event := entry.Messaging[0]
		//err = handleMessage(event.Sender.ID, event.Message.Text)
		err = handleAttachment(event.Sender.ID, event.Message.Text)
		if err != nil {
			log.Printf("Failed sending message: %s", err)
			resp.WriteHeader(400)
			resp.Write([]byte("An error occurred"))
			return
		}
	}
}

func sendRequest(data []byte) error {
	uri := "https://graph.facebook.com/v2.6/me/messages"
	uri = fmt.Sprintf("%s?access_token=%s", uri, os.Getenv("FACEBOOK_ACCESS_TOKEN"))
	log.Printf("URI: %s", uri)
	req, err := http.NewRequest(
		"POST",
		uri,
		bytes.NewBuffer(data),
	)
	if err != nil {
		log.Printf("Failed making request: %s", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed doing request: %s", err)
		return err
	}

	log.Printf("MESSAGE SENT?\n%#v", res)
	return nil
}

// Handles messages
func handleAttachment(senderId, message string) error {
	if len(message) == 0 {
		return errors.New("No message found.")
	}

	response := ResponseAttachment{
		Recipient: Recipient{
			ID: senderId,
		},
		Message: Attachment{},
	}

	elements := []Element{
		Element{
			Title:    "Check us out",
			ImageURL: "https://niceable.co/images/heart.jpg",
			Subtitle: "Fresh, organic and ethical giveaways",
			DefaultAction: DefaultAction{
				Type:                "web_url",
				URL:                 "https://niceable.co",
				WebViewHeightRation: "tall",
			},
			Buttons: []Button{
				Button{
					Type:  "web_url",
					URL:   "https://niceable.co",
					Title: "Join giveaways",
				},
			},
		},
	}

	response.Message.Attachment.Type = "template"
	response.Message.Attachment.Payload.TemplateType = "generic"
	response.Message.Attachment.Payload.Elements = elements

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal error: %s", err)
		return err
	}

	log.Printf("DATA: %s", string(data))

	return sendRequest(data)
}

// Handles messages
func handleMessage(senderId, message string) error {
	if len(message) == 0 {
		return errors.New("No message found.")
	}

	response := ResponseMessage{
		Recipient: Recipient{
			ID: senderId,
		},
		Message: Message{
			Text: "Hello",
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal error: %s", err)
		return err
	}

	return sendRequest(data)
}

// Initialize request
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handleWebhook).Methods("POST", "GET")

	port := ":8000"
	log.Printf("Server started on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
