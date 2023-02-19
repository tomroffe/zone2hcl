
test:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out

coverage:
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	open coverage/coverage.html