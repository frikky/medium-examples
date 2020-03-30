#!/bin/sh
curl -X POST -H "Content-Type: application/json" -d '{
	"get_started": {"payload": "GET_STARTED"},
  "greeting": [
    {
      "locale":"default",
      "text":"Hello {{user_first_name}}!" 
    }, {
      "locale":"en_US",
      "text":"Giveaways with a twist"
    }
  ],
	"persistent_menu": [
			{
					"locale": "default",
					"composer_input_disabled": false,
					"call_to_actions": [{
							"type": "postback",
							"title": "Join newsletter",
							"payload": "JOIN_NEWSLETTER"
					},
					{
							"type": "postback",
							"title": "Enter a giveaway",
							"payload": "GET_RAFFLES"
					},
					{
							"type": "web_url",
							"title": "Visit us",
							"url": "https://niceable.co",
							"webview_height_ratio": "full"
					}]
			}
   ]
}' "https://graph.facebook.com/v2.6/me/messenger_profile?access_token=$FACEBOOK_ACCESS_TOKEN"
