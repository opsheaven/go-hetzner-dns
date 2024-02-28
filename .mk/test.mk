.PHONY: test
test:
	go test -coverprofile=.test.out

.PHONY: cover
cover: test
	go tool cover -html=.test.out
