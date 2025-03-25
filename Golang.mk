
.PHONY: get-deps
get-deps:
	go get -v -t ./...

.PHONY: update-deps
update-deps:
	go get -u ./...
	go mod tidy

