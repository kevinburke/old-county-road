.PHONY: assets static templates

SHELL = /bin/bash

BENCHSTAT := $(GOPATH)/bin/benchstat
BUMP_VERSION := $(GOPATH)/bin/bump_version
DIFFER := $(GOPATH)/bin/differ
GO_BINDATA := $(GOPATH)/bin/go-bindata
JUSTRUN := $(GOPATH)/bin/justrun
MEGACHECK := $(GOPATH)/bin/megacheck
RELEASE := $(GOPATH)/bin/github-release

# Add files that change frequently to this list.
WATCH_TARGETS = static/style.css templates/index.html main.go

test: vet
	go test ./...

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

vet: $(MEGACHECK)
	$(MEGACHECK) --ignore='github.com/kevinburke/go-html-boilerplate/*.go:U1000' ./...
	go vet ./...

race-test: vet
	go test -race ./...

diff: $(DIFFER)
	differ $(MAKE) assets

$(BENCHSTAT):
	go get golang.org/x/perf/cmd/benchstat

bench: | $(BENCHSTAT)
	tmp=$$(mktemp); go list ./... | grep -v vendor | xargs go test -benchtime=2s -bench=. -run='^$$' > "$$tmp" 2>&1 && $(BENCHSTAT) "$$tmp"

serve:
	go install . && go-html-boilerplate

generate_cert:
	go run "$$(go env GOROOT)/src/crypto/tls/generate_cert.go" --host=localhost:7065,127.0.0.1:7065 --ecdsa-curve=P256 --ca=true

$(GO_BINDATA):
	go get -u github.com/jteeuwen/go-bindata/...

assets: | $(GO_BINDATA)
	$(GO_BINDATA) -o=assets/bindata.go --nocompress --nometadata --pkg=assets templates/... static/...

$(JUSTRUN):
	go get -u github.com/jmhodges/justrun

watch: | $(JUSTRUN)
	$(JUSTRUN) -v --delay=100ms -c 'make assets serve' $(WATCH_TARGETS)

$(BUMP_VERSION):
	go get github.com/Shyp/bump_version

$(DIFFER):
	go get github.com/kevinburke/differ

$(RELEASE):
	go get -u github.com/aktau/github-release

# Run "GITHUB_TOKEN=my-token make release version=0.x.y" to release a new version.
release: diff race-test | $(BUMP_VERSION) $(RELEASE)
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	bump_version --version=$(version) main.go
	git push origin --tags
	mkdir -p releases/$(version)
	# Change the binary names below to match your tool name
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/go-html-boilerplate-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/go-html-boilerplate-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/go-html-boilerplate-windows-amd64 .
	# Change the Github username to match your username.
	# These commands are not idempotent, so ignore failures if an upload repeats
	$(RELEASE) release --user kevinburke --repo go-html-boilerplate --tag $(version) || true
	$(RELEASE) upload --user kevinburke --repo go-html-boilerplate --tag $(version) --name go-html-boilerplate-linux-amd64 --file releases/$(version)/go-html-boilerplate-linux-amd64 || true
	$(RELEASE) upload --user kevinburke --repo go-html-boilerplate --tag $(version) --name go-html-boilerplate-darwin-amd64 --file releases/$(version)/go-html-boilerplate-darwin-amd64 || true
	$(RELEASE) upload --user kevinburke --repo go-html-boilerplate --tag $(version) --name go-html-boilerplate-windows-amd64 --file releases/$(version)/go-html-boilerplate-windows-amd64 || true
