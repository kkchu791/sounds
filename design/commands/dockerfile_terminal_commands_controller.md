### Dockerfile Terminal Commands for Controller

`$ docker buildx build -t sounds-controller .`

### command to build the broker

##### How I run my controller locally:

`$ go run ./cmd/controller`

#### Running the Docker Image for Controller (detached):

`$ docker run -d -p 9000:9000 sounds-controller`


#### Docker Compose with Broker and Controller
// using compose, it handles env variable flags and networks and does docker run

`$ docker compose build controller` #{builds the service}

`$ docker compose up broker-0` // also runs the controller because of depends_on

`$ docker compose up` // to run both broker and controller

`$ docker compose down`

`$ docker compose up --build` // you need to add build flag to get latest changes of code
