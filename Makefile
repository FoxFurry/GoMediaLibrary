tests:
	@docker-compose -f docker-compose-test.yaml up --build --abort-on-container-exit
	@docker image prune -f
	@docker-compose -f docker-compose-test.yaml down --volumes

## Run Integration Test
## Note: This command is intended to be executed within docker env
integration-tests:
	@sh -c "while ! pg_isready -d medialibrary -h postgres -p 5432 -U postgres; do echo Waiting for postgres 3s; sleep 3; done"
	@echo "Connection successful. Running tests"
	go test -v -coverprofile=./coverage.out ./...
	@echo "Tests complete. Generating code coverage"
	go tool cover -html=coverage.out -o ./coverage/coverage.html
