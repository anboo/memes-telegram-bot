TEST_CMD=$(shell type gotestsum 2>/dev/null >/dev/null && echo "gotestsum --" || echo "go test")

mocks:
	mockgen -source handler/interfaces.go -destination=handler/mocks.go  -package=handler
	mockgen -source handler/choose_age/interfaces.go -destination=handler/choose_age/mocks.go  -package=choose_age

test:
	$(TEST_CMD) -tags="testing" -v -race -cover -coverprofile=coverage.out ./...

