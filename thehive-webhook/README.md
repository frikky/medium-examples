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

## Elastic issues:
```
docker exec -u 0 -it thehive_elasticsearch_1 curl -XPUT -H "Content-Type: application/json" http://localhost:9200/_cluster/settings -d '{ "transient": { "cluster.routing.allocation.disk.threshold_enabled": false } }'; curl -XPUT -H "Content-Type: application/json" http://localhost:9200/_all/_settings -d '{"index.blocks.read_only_allow_delete": null}'
curl -XPUT http://localhost:9200/_cluster/settings -d '{ "transient": { "cluster.routing.allocation.disk.threshold_enabled": false } }' -H "Content-Type: application/json"
```
