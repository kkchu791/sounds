### Dockerfile Terminal Commands for Controller

`$ docker buildx build -t sounds-controller .`

##### How I run my controller locally:

`$ go run ./cmd/controller`

#### Running the Docker Image for Controller (detached):

`$ docker run -d -p 9000:9000 sounds-controller`
