#!/bin/sh
curl -X POST -H "Content-Type: application/json" -d '{
  "greeting": [
    {
      "locale":"default",
      "text":"Hello {{user_first_name}}!" 
    }, {
      "locale":"en_US",
      "text":"Giveaways with a twist"
    }
  ]
}' "https://graph.facebook.com/v2.6/me/messenger_profile?access_token=$FACEBOOK_ACCESS_TOKEN"

curl -X POST -H "Content-Type: application/json" -d '{
	  "get_started": {"payload": "GET_STARTED"}
}' "https://graph.facebook.com/v2.6/me/messenger_profile?access_token=$FACEBOOK_ACCESS_TOKEN"

