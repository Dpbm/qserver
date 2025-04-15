TEST_DIR = ./test
USAGE_TESTS = ./usage-tests.sh
RUN_SCRIPT = ./run-server-locally-for-tests.sh
SWAGGER_SCRIPT = ./generate-swagger.sh

install:
	go mod download 

clear-docker:
	docker compose -f ../../dev-compose.yml down -v

go-test:
	go test $(TEST_DIR)

test-docker: go-test clear-docker
	docker compose -f ../../dev-compose.yml up -d --build jobs-server api
	cd $(TEST_DIR) && chmod +x $(USAGE_TESTS) && HOST=0.0.0.0 $(USAGE_TESTS)

test-local: go-test clear-docker
	docker compose -f ../../dev-compose.yml up -d --build jobs-server 
	echo "Waiting docker to start" && sleep 10
	cd $(TEST_DIR) && chmod +x $(USAGE_TESTS) && HOST=localhost $(USAGE_TESTS)

run:
	chmod +x $(RUN_SCRIPT) && $(RUN_SCRIPT)

swagger:
	chmod +x $(SWAGGER_SCRIPT) && $(SWAGGER_SCRIPT)

lint:
	golangci-lint run