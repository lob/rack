.PHONY: all build clean release

all: build

build:
	go install ./...

clean:
	rm -f pkg/convox-*

release: clean
ifndef VERSION
	$(error VERSION must be specified)
endif
	go get -u github.com/karalabe/xgo
	$(GOPATH)/bin/xgo -branch $(shell git rev-parse HEAD) -out pkg/convox -targets 'darwin/amd64,linux/amd64,windows/amd64' -ldflags "-X main.version=$(VERSION)" .
	aws-creds -p core
	aws s3 cp pkg/convox-darwin-10.6-amd64 s3://public.lob.com/engineering/convox/darwin/convox --acl public-read --profile core
	aws s3 cp pkg/convox-linux-amd64 s3://public.lob.com/engineering/convox/linux/convox --acl public-read --profile core
	aws s3 cp pkg/convox-windows-4.0-amd64.exe s3://public.lob.com/engineering/convox/windows/convox.exe --acl public-read --profile core
