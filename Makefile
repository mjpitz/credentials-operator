build-deps:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups
	GO111MODULE=off go get -u github.com/mitchellh/gox

deps:
	go mod download
	go mod verify

fmt:
	go-groups -w .
	gofmt -s -w .

test:
	go vet ./...
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

deploy:
	mkdir -p bin
	gox -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

docker:
	docker build -t mjpitz/credentials-operator:latest .
