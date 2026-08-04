### Setup Nodes and Create a topic
Run the controller with:
Terminal 1:
`go run ./cmd/controller`

Run the leader with:
Terminal 2 (leader):
` BROKER_ID=broker-0 BROKER_ADDR=localhost:9001 BROKER_PORT=9001 CONTROLLER_ADDR=localhost:9000 go run ./cmd/broker`

Run the follower with:
Terminal 2 (follower):

` BROKER_ID=broker-1 BROKER_ADDR=localhost:9002 BROKER_PORT=9002 CONTROLLER_ADDR=localhost:9000 go run ./cmd/broker`


Run the follower with:
Terminal 3 (follower):

 `BROKER_ID=broker-2 BROKER_ADDR=localhost:9003 BROKER_PORT=9003 CONTROLLER_ADDR=localhost:9000 go run ./cmd/broker`

After brokers are up run admin client:
Terminal 4 (admin client):

`export BOOTSTRAP_SERVERS="http://localhost:9001,http://localhost:9002"`

`go run ./cmd/admin -- create-topic --topic Tokyo-Sounds --partitions 6 --replication-factor 3`

### Sending a message to the leader
1. Append to the leader:
`curl -X POST localhost:9001/append -H "Content-Type: application/json" -d '{"sound":"rain","key":"outdoor"}'`

2. Wait a second, then read from the follower:
`curl "localhost:9002/read?offset=0"`
