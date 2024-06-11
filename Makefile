container:
	@docker-compose build test
travis:
	@go get -v -d -t
test:
	@docker-compose run --rm test
coverage:
	@docker-compose run --rm test sh -c 'go test -coverprofile=coverage.out ./... && go tool cover -html coverage.out -o coverage.html'
	@open coverage.html
clean:
	@docker image rm go-shopify
	@rm -f coverage.html coverage.out

.DEFAULT_GOAL := container
