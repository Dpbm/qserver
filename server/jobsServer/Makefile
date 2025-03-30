TEST_SCRIPT = ./test.sh
RUN_SCRIPT = ./run-server-for-local-tests.sh
PROTO_DIR = ./proto
GEN_PROTO_SCRIPT = ./generate-proto-code.sh


test:
	chmod +x $(TEST_SCRIPT) && $(TEST_SCRIPT)

run:
	chmod +x $(RUN_SCRIPT) && $(RUN_SCRIPT)

proto:
	cd $(PROTO_DIR) && chmod +x $(GEN_PROTO_SCRIPT) && $(GEN_PROTO_SCRIPT)