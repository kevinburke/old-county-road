.PHONY: assets static templates

SHELL = /bin/bash

GO_BINDATA := $(shell command -v go-bindata)
JUSTRUN := $(shell command -v justrun)
STATICCHECK := $(shell command -v staticcheck)

# Add files that change frequently to this list.
WATCH_TARGETS = static/style.css templates/index.html main.go

vet:
ifndef STATICCHECK
	go get -u honnef.co/go/tools/cmd/staticcheck
endif
	staticcheck ./...
	go vet ./...

test: vet
	go test ./...

race-test: vet
	go test -race ./...

serve:
	go install . && go-html-boilerplate

generate_cert:
	go run "$$(go env GOROOT)/src/crypto/tls/generate_cert.go" --host=localhost:7065,127.0.0.1:7065 --ecdsa-curve=P256 --ca=true

assets:
ifndef GO_BINDATA
	go get -u github.com/jteeuwen/go-bindata/...
endif
	go-bindata -o=assets/bindata.go --nometadata --pkg=assets templates/... static/...

watch:
ifndef JUSTRUN
	go get -u github.com/jmhodges/justrun
endif
	justrun -v --delay=100ms -c 'make assets serve' $(WATCH_TARGETS)
