# Design: Controller→Broker LeaderAndIsr Propagation

### Problem/Context
Brokers need to learn partition leadership and ISR changes decided by the controller. Currently no mechanism exists to propagate this state, so replicas can't correctly reject/accept produce requests.

### Proposed Solution
Controller pushes LeaderAndIsrRequest to affected brokers over HTTP whenever `PartitionStateMachine` transitions. Alternative considered: brokers poll controller metadata on an interval — rejected due to...

### Architecture / Technical Details

New LeaderAndIsrRequest/Response structs service package
Controller sends are fire-and-forget
Broker's handler updates local PartitionReplica state

### Appendix
When does ISRleader fire?
