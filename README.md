
Run the leader with:
Terminal 1 (leader):
`IS_LEADER=true BROKER_ID=broker-0 BROKER_PORT=5001 go run ./cmd/broker`

Run the follower with:
Terminal 2 (follower):
`IS_LEADER=false BROKER_ID=broker-1 BROKER_PORT=5002 LEADER_ADDR=localhost:5001 go run ./cmd/broker`

1. Append to the leader:
`curl -X POST localhost:5001/append -H "Content-Type: application/json" -d '{"sound":"rain","key":"outdoor"}'`

2. Wait a second, then read from the follower:
`curl "localhost:5002/read?offset=0"`