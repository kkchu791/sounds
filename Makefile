.PHONY: cluster stop controller brokers create-topic logs clean

PID_DIR := .pids
LOG_DIR := .logs

CONTROLLER_ADDR := localhost:9000
BOOTSTRAP_SERVERS := http://localhost:9001,http://localhost:9002

.PHONY: dirs
dirs:
	mkdir -p $(PID_DIR) $(LOG_DIR)

controller: dirs
	go run ./cmd/controller > $(LOG_DIR)/controller.log 2>&1 & echo $$! > $(PID_DIR)/controller.pid
	@echo "controller starting (logs: $(LOG_DIR)/controller.log)"

brokers: dirs
	BROKER_ID=broker-0 BROKER_ADDR=localhost:9001 BROKER_PORT=9001 CONTROLLER_ADDR=$(CONTROLLER_ADDR) \
		go run ./cmd/broker > $(LOG_DIR)/broker-0.log 2>&1 & echo $$! > $(PID_DIR)/broker-0.pid
	BROKER_ID=broker-1 BROKER_ADDR=localhost:9002 BROKER_PORT=9002 CONTROLLER_ADDR=$(CONTROLLER_ADDR) \
		go run ./cmd/broker > $(LOG_DIR)/broker-1.log 2>&1 & echo $$! > $(PID_DIR)/broker-1.pid
	BROKER_ID=broker-2 BROKER_ADDR=localhost:9003 BROKER_PORT=9003 CONTROLLER_ADDR=$(CONTROLLER_ADDR) \
		go run ./cmd/broker > $(LOG_DIR)/broker-2.log 2>&1 & echo $$! > $(PID_DIR)/broker-2.pid
	@echo "3 brokers starting (logs: $(LOG_DIR)/broker-*.log)"

create-topic:
	BOOTSTRAP_SERVERS=$(BOOTSTRAP_SERVERS) \
		go run ./cmd/admin -- create-topic --topic Tokyo-Sounds --partitions 6 --replication-factor 3

cluster: controller
	sleep 1
	$(MAKE) brokers
	sleep 2
	$(MAKE) create-topic
	@echo "cluster is up. tail logs with: make logs"

logs:
	tail -f $(LOG_DIR)/*.log

stop:
	-@for f in $(PID_DIR)/*.pid; do \
		[ -f "$$f" ] && kill $$(cat "$$f") 2>/dev/null; \
	done
	@rm -rf $(PID_DIR)
	@echo "cluster stopped"

clean: stop
	rm -rf $(LOG_DIR)