TEST_CMD=$(shell type gotestsum 2>/dev/null >/dev/null && echo "gotestsum --" || echo "go test")

mocks:
	mockgen -source handler/interfaces.go -destination=handler/mocks.go  -package=handler
	mockgen -source handler/choose_age/interfaces.go -destination=handler/choose_age/mocks.go  -package=choose_age
	mockgen -source handler/choose_sex/interfaces.go -destination=handler/choose_sex/mocks.go  -package=choose_sex
	mockgen -source handler/send_mem/interfaces.go -destination=handler/send_mem/mocks.go  -package=send_mem
	mockgen -source handler/vote/interfaces.go -destination=handler/vote/mocks.go  -package=vote
	mockgen -source handler/welcome/interfaces.go -destination=handler/welcome/mocks.go  -package=welcome
	mockgen -source handler/config_source/interfaces.go -destination=handler/config_source/mocks.go  -package=config_source

test:
	$(TEST_CMD) -tags="testing" -v -race -cover -coverprofile=coverage.out ./...

