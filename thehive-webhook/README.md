# TheHive conf
This is related to part 4 of the Shuffle blogpost series.

## Get started with Webhook testing
1. Open application.conf
2. Edit the line url = with your url as such:
```
webhooks {
  myLocalWebHook {
		url = "YOUR URL HERE"
  }
}
```
3. Run docker-compose
```
docker-compose up
```
4. Go to http://localhost:9000 and set up users
5. Create an example case
6. Check Shuffle whether the webhook executed
7. ???
8. Profit!
