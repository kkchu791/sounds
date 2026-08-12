# Design: Controller→Broker LeaderAndIsr Propagation

### Problem/Context
Brokers need to learn partition leadership and ISR changes from the controller. Currently no mechanism exists to propagate this state, so replicas can't correctly reject/accept produce requests.

### Proposed Solution
Brokers periodically pull MetadataRequest and LeaderAndIsrRequest from controller. 

Alternatives considered: Controllers pushes Metadata and LeaderAndIsrRequest when Controller's PartitionStateMachine changes - rejected because with the Push Approach, there is inconsistency between brokers when there are network crashes, retries eventually stop, and hard to manage timing of multiple retries, which leads to data being written in the wrong order. (KIP-500: Replace ZooKeeper with a Self-Managed Metadata Quorum)

### Architecture / Technical Details

Brokerss fetch from the Controllerss
Broker's handler updates local PartitionReplica state

### Appendix

