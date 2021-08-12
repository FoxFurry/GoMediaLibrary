tests:
	@docker-compose -f docker-compose-test.yaml up --build --abort-on-container-exit
	@docker-compose -f docker-compose-test.yaml down --volumes

## Run Integration Test
## Note: This command is intended to be executed within docker env
integration-tests:
	@sh -c "while ! pg_isready -d medialibrary -h postgres -p 5432 -U postgres; do echo Waiting for postgres 3s; sleep 3; done"
	go test -v ./...