##### create KAFKA_CLUSTER_ID
```KAFKA_CLUSTER_ID="$(bin/kafka-storage.sh random-uuid)"```

##### Generate directory ids
```
CONTROLLER_1_UUID="$(bin/kafka-storage.sh random-uuid)"
CONTROLLER_2_UUID="$(bin/kafka-storage.sh random-uuid)"
CONTROLLER_2_UUID="$(bin/kafka-storage.sh random-uuid)"
```

##### server.property edits

create 3 seperate server.properties files and change these fields:

```
node.id={broker #}
controller.quorum.bootstrap.servers=localhost:{9093},localhost:{9095},localhost:{9097}
listeners=PLAINTEXT://:{9092},CONTROLLER://:{9093}
advertised.listeners=PLAINTEXT://localhost:{9092},CONTROLLER://localhost:{9093} //change ports
log.dirs=/tmp/kraft-logs-{broker #} // where to store the log files
```


##### format each brokers storage - provide the server its configs from server.properties
{change server-{num} for each server}

```bin/kafka-storage.sh format --cluster-id $KAFKA_CLUSTER_ID --initial-controllers "1@localhost:9093:${CONTROLLER_1_UUID},2@localhost:9095:${CONTROLLER_2_UUID},3@localhost:9097:${CONTROLLER_3_UUID}" -c config/server-1.properties```

##### start brokers
`bin/kafka-server-start.sh config/server-1.properties`


##### create topic with partition count and replication factor count
```bin/kafka-topics.sh --create --topic Sounds --partitions 6 --replication-factor 3 --bootstrap-server localhost:9092```

##### describe the topic
```bin/kafka-topics.sh --describe --topic Sounds --bootstrap-server localhost:9092```
