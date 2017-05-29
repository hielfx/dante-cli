.SILENT :
.PHONY : dante-cli clean fmt

TAG:=`git describe --abbrev=0 --tags`
LDFLAGS:=-X main.buildVersion=$(TAG)

all: dante-cli

deps:
	go get -u github.com/FiloSottile/gvt
	gvt list

dante-cli:
	echo "Building dante-cli"
	go install -ldflags "$(LDFLAGS)"

dist-clean:
	rm -rf dist
	rm -f dante-cli-alpine-linux-*.tar.gz
	rm -f dante-cli-linux-*.tar.gz

dist: deps dist-clean
	mkdir -p dist/alpine-linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -tags netgo -installsuffix netgo -o dist/alpine-linux/amd64/dante-cli
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux/amd64/dante-cli
	mkdir -p dist/linux/armel && GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(LDFLAGS)" -o dist/linux/armel/dante-cli
	mkdir -p dist/linux/armhf && GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "$(LDFLAGS)" -o dist/linux/armhf/dante-cli

release: dist
	tar -cvzf dist/dante-cli-alpine-linux-amd64-$(TAG).tar.gz -C dist/alpine-linux/amd64 dante-cli
	tar -cvzf dist/dante-cli-linux-amd64-$(TAG).tar.gz -C dist/linux/amd64 dante-cli
	tar -cvzf dist/dante-cli-linux-armel-$(TAG).tar.gz -C dist/linux/armel dante-cli
	tar -cvzf dist/dante-cli-linux-armhf-$(TAG).tar.gz -C dist/linux/armhf dante-cli
