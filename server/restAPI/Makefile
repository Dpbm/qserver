TEST_DIR = ./test
USAGE_TESTS = ./usage-tests.sh
RUN_SCRIPT = ./run-server-for-local-tests.sh
SWAGGER_SCRIPT = ./generate-swagger.sh

test-docker:
	cd $(TEST_DIR) && chmod +x HOST=0.0.0.0 $(USAGE_TESTS) && $(USAGE_TESTS)

test-local:
	cd $(TEST_DIR) && chmod +x HOST=localhost $(USAGE_TESTS) && $(USAGE_TESTS)

run:
	chmod +x $(RUN_SCRIPT) && $(RUN_SCRIPT)

swagger:
	chmod +x $(SWAGGER_SCRIPT) && $(SWAGGER_SCRIPT)