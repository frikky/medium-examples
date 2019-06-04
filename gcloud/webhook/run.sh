docker stop webhook
docker rm webhook
docker rmi webhook

PORT=5001

docker build . -t webhook
docker run -d \
	-e "HOOKPORT=$PORT" \
	-p $PORT:$PORT \
	--name webhook \
	-h webhook \
	--restart always \
	webhook

docker logs -f webhook
