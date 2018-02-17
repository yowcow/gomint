all: dep

dep:
	dep ensure -v

test:
	go test -v ./...
