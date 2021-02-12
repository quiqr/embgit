# Go parameters
GO111MODULE=on
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=embgit
BINARY_UNIX=$(BINARY_NAME)_unix
POPPYAPPDIR=~/cPoppyGo/poppygo-eletron-app

all: test build
build:
	GO111MODULE=$(GO111MODULE) $(GOBUILD) -o $(BINARY_NAME) -v
buildx:
	gox -osarch="linux/amd64" ./
	gox -osarch="windows/amd64" ./
	gox -osarch="darwin/amd64" ./

cptopoppygopp:
	cp embgit_darwin_amd64 $(POPPYAPPDIR)/resources/mac/embgit
	cp embgit_windows_amd64.exe $(POPPYAPPDIR)/resources/win/embgit.exe
	cp embgit_linux_amd64 $(POPPYAPPDIR)/resources/linux/embgit

test:
	GO111MODULE=$(GO111MODULE) $(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)*
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/urfave/cli/v2

release:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is not defined)
endif
	git commit -am 'Update version to $(version)'
	git tag -a $(version) -m '$(version)'
	git push origin $(version)
	goreleaser --rm-dist

# Cross compilation
#build-linux:
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
#docker-build:
#	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v

inttest:
#				@rm -Rfv /tmp/test
#				./embgit clone -i ~/.ssh/id_rsa-annemarie-vega git@gitlab.lingewoud.net:servers/scio-site-monitor.git /tmp/test
#				@rm -Rfv /tmp/test
#				./embgit clone -i ~/.ssh/id_rsa-annemarie-vega git@github.com:mipmip/linny.vim.git /tmp/test
#				@rm -Rfv /tmp/test
#				./embgit clone git@github.com:mipmip/linny.vim.git /tmp/test
	@rm -Rfv /tmp/test
	./embgit clone -i ~/.ssh/id_rsa-annemarie-vega git@gitlab.lingewoud.net:Sandbox/testembgit.git /tmp/test
	echo ".\n" >> /tmp/test/test
	mkdir /tmp/test/sub
	touch /tmp/test/sub/yoehoe
	./embgit alladd /tmp/test
	./embgit commit -n "Pim Snel" -e "pim@lingewoud.nl" -m "a message" /tmp/test
	./embgit push -i ~/.ssh/id_rsa-annemarie-vega /tmp/test
