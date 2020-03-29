# Go parameters
GO111MODULE=on
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=embgit
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
				GO111MODULE=$(GO111MODULE) $(GOBUILD) -o $(BINARY_NAME) -v
test:
				GO111MODULE=$(GO111MODULE) $(GOTEST) -v ./...
clean:
				$(GOCLEAN)
				rm -f $(BINARY_NAME)
				rm -f $(BINARY_UNIX)
run:
				$(GOBUILD) -o $(BINARY_NAME) -v ./...
				./$(BINARY_NAME)
deps:
				$(GOGET) github.com/urfave/cli/v2

# Cross compilation
build-linux:
				CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
				docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v

inttest:
				@rm -Rfv /tmp/test
				./embgit clone -i ~/.ssh/id_rsa-annemarie-vega git@gitlab.lingewoud.net:servers/scio-site-monitor.git /tmp/test
				@rm -Rfv /tmp/test
				./embgit clone -i ~/.ssh/id_rsa-annemarie-vega git@github.com:mipmip/linny.vim.git /tmp/test
				@rm -Rfv /tmp/test
				./embgit clone git@github.com:mipmip/linny.vim.git /tmp/test
