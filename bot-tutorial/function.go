package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var newsLetterIds = []string{}
var raffleIds = map[string][]string{}

// The message we get from Messenger
type InputMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Postback struct {
				Title   string `json:"title"`
				Payload string `json:"payload"`
			} `json:"postback"`
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
						Email []struct {
							Confidence float64 `json:"confidence"`
							Value      string  `json:"value"`
						} `json:"email"`
					} `json:"entities"`
					DetectedLocales []struct {
						Locale     string  `json:"locale"`
						Confidence float64 `json:"confidence"`
					} `json:"detected_locales"`
				} `json:"nlp"`
				QuickReply struct {
					Payload string `json:"payload"`
				} `json:"quick_reply"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

// The recipient of our message
type Recipient struct {
	ID string `json:"id"`
}

type QuickReply struct {
	ContentType string `json:"content_type,omitempty"`
	Title       string `json:"title,omitempty"`
	Payload     string `json:"payload,omitempty"`
	ImageUrl    string `json:"image_url,omitempty"`
}

// The message to send it its basic
type Message struct {
	Text         string       `json:"text,omitempty"`
	Mid          string       `json:"mid,omitempty"`
	QuickReplies []QuickReply `json:"quick_replies,omitempty"`
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
	Recipient   Recipient  `json:"recipient"`
	MessageType string     `json:"message_type,omitempty"`
	Message     Attachment `json:"message,omitempty"`
}

// Full response
type ResponseMessage struct {
	Recipient     Recipient `json:"recipient"`
	MessagingType string    `json:"messaging_type,omitempty"`
	Message       Message   `json:"message,omitempty"`
}

func HandleMessenger(resp http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		secretKey := os.Getenv("FACEBOOK_SECRET_KEY")
		u, err := url.Parse(request.RequestURI)
		if err != nil {
			log.Printf("Failed parsing url: %s", err)
			resp.WriteHeader(400)
			resp.Write([]byte(fmt.Sprintf("An error occurred: %s", err)))
			return
		}

		values, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			log.Printf("Failed parsing values url: %s", err)
			resp.WriteHeader(400)
			resp.Write([]byte(fmt.Sprintf("An error occurred: %s", err)))
			return
		}

		token := values.Get("hub.verify_token")
		if token == secretKey {
			resp.WriteHeader(200)
			resp.Write([]byte(values.Get("hub.challenge")))
			return
		}

		log.Printf("VALUES: %#v", values)

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
	for _, entry := range message.Entry {
		if len(entry.Messaging) == 0 {
			log.Printf("No messages")
			resp.WriteHeader(400)
			resp.Write([]byte("An error occurred"))
			return
		}

		event := entry.Messaging[0]

		// Handle the email with the MID we had earlier
		if len(event.Message.Nlp.Entities.Email) > 0 {
			if event.Message.Nlp.Entities.Email[0].Confidence > 0.99 {
				err = handleEmail(event.Message.Nlp.Entities.Email[0].Value, event.Sender.ID)
				if err != nil {
					log.Printf("Failed handling email: %s", err)
					resp.WriteHeader(400)
					resp.Write([]byte("An error occurred"))
					return
				}
				continue
			}
		}

		if len(event.Postback.Payload) > 0 {
			err = handlePostback(event.Sender.ID, event.Postback.Title, event.Postback.Payload)
			if err != nil {
				log.Printf("Failed payload(1): %s", err)
				resp.WriteHeader(400)
				resp.Write([]byte("An error occurred"))
				return
			}
		} else if len(event.Message.QuickReply.Payload) > 0 {
			err = handlePostback(event.Sender.ID, event.Message.Text, event.Message.QuickReply.Payload)
			if err != nil {
				log.Printf("Failed payload(2): %s", err)
				resp.WriteHeader(400)
				resp.Write([]byte("An error occurred"))
				return
			}
		} else if len(event.Message.Text) > 0 {
			err = handleText(event.Sender.ID, event.Message.Text)
			if err != nil {
				log.Printf("Failed payload(3 - text): %s", err)
				resp.WriteHeader(400)
				resp.Write([]byte("An error occurred"))
				return
			}
		}
	}
}

func handleEmail(email, senderId string) error {
	// Newsletter
	for index, item := range newsLetterIds {
		if item == senderId {
			// Reply
			// Sign up for newsletter
			// Remove item
			log.Printf("FOUND!")
			response := ResponseMessage{
				Recipient: Recipient{
					ID: senderId,
				},
				Message: Message{
					Text: fmt.Sprintf("%s is now signed up to receive our newsletter! :)", email),
				},
			}

			data, err := json.Marshal(response)
			if err != nil {
				log.Printf("Marshal error: %s", err)
				return err
			}

			err = sendRequest(data)
			if err != nil {
				return err
			}

			newsLetterIds = append(newsLetterIds[:index], newsLetterIds[index+1:]...)
			break
		}
	}

	return nil
}

func handleJoinNewsletterMessenger(senderId string) error {
	response := ResponseMessage{
		Recipient: Recipient{
			ID: senderId,
		},
		Message: Message{
			Text: "You are now signed up to receive our newsletter through messenger! :)",
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal error: %s", err)
		return err
	}

	err = sendRequest(data)
	if err != nil {
		return err
	}

	return nil
}

func handleJoinNewsletter(senderId string) error {
	response := ResponseMessage{
		Recipient: Recipient{
			ID: senderId,
		},
		MessagingType: "RESPONSE",
		Message: Message{
			Text: "How do you want to sign up?",
			QuickReplies: []QuickReply{
				QuickReply{
					ContentType: "text",
					Title:       "Messenger",
					Payload:     "ENTER_NEWSLETTER_MESSENGER",
				},
				QuickReply{
					ContentType: "user_email",
				},
			},
		},
	}

	// TMP
	log.Printf("Adding %s to newsletter handler!", senderId)
	newsLetterIds = append(newsLetterIds, senderId)

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal error: %s", err)
		return err
	}

	err = sendRequest(data)
	if err != nil {
		return err
	}

	return nil
}

func handleGetStarted(senderId string) error {
	response := ResponseMessage{
		Recipient: Recipient{
			ID: senderId,
		},
		MessagingType: "RESPONSE",
		Message: Message{
			Text: "Hey, what can we help you with?",
			QuickReplies: []QuickReply{
				QuickReply{
					ContentType: "text",
					Title:       "Enter a giveaway",
					Payload:     "GET_RAFFLES",
				},
				QuickReply{
					ContentType: "text",
					Title:       "Join newsletter",
					Payload:     "JOIN_NEWSLETTER",
				},
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Marshal error: %s", err)
		return err
	}

	err = sendRequest(data)
	if err != nil {
		return err
	}

	return nil
}

func handleGetRaffles(senderId string) error {
	response := ResponseAttachment{
		Recipient: Recipient{
			ID: senderId,
		},
		Message: Attachment{},
	}

	elements := []Element{
		Element{
			Title:    "Win a Huge Box of Japanese Snacks!",
			Subtitle: "No one does flavors quite like Japan. Enter this raffle to win a HUGE BOX OF JAPANESE SNACKS! Hard to find anywhere else!",
			ImageURL: "https://firebasestorage.googleapis.com/v0/b/kaechan.appspot.com/o/raffles%2F229554bf-53ba-4c5e-b9c3-2e8211cb6e60%2FWin-a-Huge-Box-of-Japanese-Snacks!.png?alt=media&token=4a4e9c7d-a077-40f3-b8f0-557af1bd0b53",
			DefaultAction: DefaultAction{
				Type:                "web_url",
				URL:                 "https://niceable.io/raffles/229554bf-53ba-4c5e-b9c3-2e8211cb6e60",
				WebViewHeightRation: "tall",
			},
			Buttons: []Button{
				Button{
					Type:  "web_url",
					URL:   "https://niceable.io/raffles/229554bf-53ba-4c5e-b9c3-2e8211cb6e60",
					Title: "Learn more",
				},
			},
		},
		Element{
			Title:    "Win adorable alpaca prizes!",
			Subtitle: "Alpaca wool is extraordinarily soft and warm. Enter this raffle to win an alpaca wool hat of your choice and a cute alpaca souvenir from Alpakafarm!",
			ImageURL: "https://firebasestorage.googleapis.com/v0/b/kaechan.appspot.com/o/raffles%2F42497cea-de22-4fbc-8ec5-40af98954403%2FWin-amazing-alpaca-prizes!.png?alt=media&token=790d8b72-0b65-4a6b-850e-14e12461920a",
			DefaultAction: DefaultAction{
				Type:                "web_url",
				URL:                 "https://niceable.io/raffles/42497cea-de22-4fbc-8ec5-40af98954403",
				WebViewHeightRation: "tall",
			},
			Buttons: []Button{
				Button{
					Type:  "web_url",
					URL:   "https://niceable.io/raffles/42497cea-de22-4fbc-8ec5-40af98954403",
					Title: "Learn more",
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

	err = sendRequest(data)
	if err != nil {
		return err
	}

	return nil
}

// Handles messages
func handleText(senderId, message string) error {
	if len(message) == 0 {
		return nil
	}

	log.Printf("INSIDE POSTBACK. Title: message", message)
	msg := strings.ToUpper(message)
	if strings.Contains(msg, "HELP") || strings.Contains(msg, "GET STARTED") {
		return handleGetStarted(senderId)
	}

	return nil
}

// Handles postback messages
func handlePostback(senderId, title, postback string) error {
	if len(postback) == 0 {
		return nil
	}

	log.Printf("INSIDE POSTBACK. Title: %s, postback: %s", title, postback)
	switch message := postback; message {
	case "GET_STARTED":
		return handleGetStarted(senderId)
	case "GET_RAFFLES":
		return handleGetRaffles(senderId)
	case "JOIN_NEWSLETTER":
		return handleJoinNewsletter(senderId)
	case "ENTER_NEWSLETTER_MESSENGER":
		return handleJoinNewsletterMessenger(senderId)
	default:
		log.Printf("Can't handle payload %s", postback)
		return nil
	}

	return nil
}

func sendRequest(data []byte) error {
	uri := "https://graph.facebook.com/v2.6/me/messages"
	uri = fmt.Sprintf("%s?access_token=%s", uri, os.Getenv("FACEBOOK_ACCESS_TOKEN"))
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

	if res.StatusCode != 200 {
		log.Printf("MESSAGE: %#v", res)
	}

	return nil
}

// Initialize request
//func main() {
//	router := mux.NewRouter()
//	router.HandleFunc("/", handleWebhook).Methods("POST", "GET")
//
//	port := ":8000"
//	log.Printf("Server started on %s", port)
//	log.Fatal(http.ListenAndServe(port, router))
//}
