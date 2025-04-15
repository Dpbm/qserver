TEST_SCRIPT = ./test.sh
RUN_SCRIPT = ./run-server-for-local-tests.sh
PROTO_DIR = ./proto
GEN_PROTO_SCRIPT = ./generate-proto-code.sh


install:
	go mod download 

test-docker: run-api
	export HOST="0.0.0.0"
	chmod +x $(TEST_SCRIPT) && $(TEST_SCRIPT)

test-local: run-api
	export HOST="local"
	chmod +x $(TEST_SCRIPT) && $(TEST_SCRIPT)

run-api:
	docker compose -f ../../dev-compose.yml up -d --build api

run:
	chmod +x $(RUN_SCRIPT) && $(RUN_SCRIPT)

proto:
	cd $(PROTO_DIR) && chmod +x $(GEN_PROTO_SCRIPT) && $(GEN_PROTO_SCRIPT)

lint:
	golangci-lint run