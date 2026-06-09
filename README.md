

1. Append to the leader:
`curl -X POST localhost:5001/append -H "Content-Type: application/json" -d '{"sound":"rain","key":"outdoor"}'`

2. Wait a second, then read from the follower:
`curl "localhost:5002/read?offset=0"`