@default: help

App := 'DrawlScan'
Version := `grep '^const VERSION = ' cmd/main/version.go | sed "s/^VERSION = \"\(.*\)\"/\1/g"`

# show help message
@help:
    echo "Build tool for {{ App }} {{ Version }} with Just"
    echo "Usage: just <recipe>"
    echo ""
    just --list

build: test
    go build -o drawlscan cmd/main/drawlscan.go

test:
    go test -covermode=count -coverprofile=coverage.out ./...