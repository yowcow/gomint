all: dep

dep:
	which dep || go get github.com/golang/dep/cmd/dep
	dep ensure -v

test:
	go test -v ./...
