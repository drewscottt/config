.PHONY: test
test:
	go test $(go list ./... | grep -v /cmd/)

.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out $(go list ./... | grep -v /cmd/)
	go tool cover -html=coverage.out

