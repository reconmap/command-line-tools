
.PHONY: get-deps
get-deps:
	go get -v -t -d ./...

.PHONY: update-deps
update-deps:
	go get -u ./...
	go mod tidy

