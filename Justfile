# show help message
@default: help

App := 'DrawlScan'
Version := `grep '^const VERSION = ' cmd/main/version.go | sed "s/^VERSION = \"\(.*\)\"/\1/g"`

# show help message
@help:
    echo "Build tool for {{ App }} {{ Version }} with Just"
    echo "Usage: just <recipe>"
    echo ""
    just --list

# bulid the applicaion with running tests
build: test
    go build -o drawlscan cmd/main/drawlscan.go cmd/main/version.go

# run tests and generate the coverage report
test:
    go test -covermode=count -coverprofile=coverage.out ./...

# clean up build artifacts
clean:
    go clean
    rm -f coverage.out drawlscan build

# update the version if the new version is provided
update_version new_version = "":
    if [ "{{ new_version }}" != "" ]; then \
        sed 's/$VERSION/{{ new_version }}/g' .template/README.md > README.md; \
        sed 's/$VERSION/{{ new_version }}/g' .template/version.go > cmd/main/version.go; \
    fi

# build DrawlScan for all platforms
make_distribution_files:
    for os in "linux" "windows" "darwin"; do \
        for arch in "amd64" "arm64"; do \
            mkdir -p dist/DrawlScan-$os-$arch; \
            if [ "$os" = "linux" ] && [ "$arch" = "amd64" ]; then \
                env GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o dist/DrawlScan-$os-$arch/DrawlScan cmd/main/drawlscan.go cmd/main/version.go; \
            else \
                env GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o dist/DrawlScan-$os-$arch/DrawlScan cmd/main/drawlscan.go cmd/main/version.go; \
            fi; \
            cp README.md LICENSE dist/DrawlScan-$os-$arch; \
            tar cvfz dist/DrawlScan-$os-$arch.tar.gz -C dist DrawlScan-$os-$arch; \
        done; \
    done

# upload assets to the GitHub release page
upload_assets tag:
    gh release upload --repo nagayon-935/{{ App }} {{ tag }} dist/*.tar.gz

# build DrawlScan for all platforms (OSごとに実行)
build_release_binaries:
  if [ "$OS" = "Linux" ]; then
    for arch in "amd64" "arm64"; do
      mkdir -p dist/DrawlScan-linux-$arch
      env GOOS=linux GOARCH=$arch CGO_ENABLED=1 go build -o dist/DrawlScan-linux-$arch/DrawlScan cmd/main/drawlscan.go cmd/main/version.go
      cp README.md LICENSE dist/DrawlScan-linux-$arch
      tar czf dist/DrawlScan-linux-$arch.tar.gz -C dist DrawlScan-linux-$arch
    done
  elif [ "$OS" = "Darwin" ]; then
    for arch in "amd64" "arm64"; do
      mkdir -p dist/DrawlScan-darwin-$arch
      env GOOS=darwin GOARCH=$arch CGO_ENABLED=1 go build -o dist/DrawlScan-darwin-$arch/DrawlScan cmd/main/drawlscan.go cmd/main/version.go
      cp README.md LICENSE dist/DrawlScan-darwin-$arch
      tar czf dist/DrawlScan-darwin-$arch.tar.gz -C dist DrawlScan-darwin-$arch
    done
  elif [ "$OS" = "Windows_NT" ]; then
    for arch in "amd64" "arm64"; do
      mkdir -p dist/DrawlScan-windows-$arch
      env GOOS=windows GOARCH=$arch CGO_ENABLED=1 go build -o dist/DrawlScan-windows-$arch/DrawlScan.exe cmd/main/drawlscan.go cmd/main/version.go
      cp README.md LICENSE dist/DrawlScan-windows-$arch
      tar czf dist/DrawlScan-windows-$arch.tar.gz -C dist DrawlScan-windows-$arch
    done
  fi