
To run for example on port 8080
```
sudo docker build -t post-back .
sudo docker container run -e BATCH_SIZE=5 -e BATCH_SECOND_INTERVAL=5 -e POST_ENDPOINT_URL=http://byc1u0wq5v3hkle4.b.requestbin.net -e GIN_MODE=release -e PORT=8080 -p 8080:8080 post-back
```