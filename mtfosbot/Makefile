
.PHONY: clean build

build:
	GO111MODULE=off go get -u github.com/otakukaze/go-bindata/...
	go-bindata -pkg schema -ignore .git -o module/schema/static.go schema/
	GOOS=linux go build -o mtfosbot -ldflags "-s -w" .

clean:
	rm -rf mtfosbot && go clean
