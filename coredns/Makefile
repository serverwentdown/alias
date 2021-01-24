GO:=go
NAME:=coredns

.PHONY: build
build:
	$(RM) -rf build && mkdir build
	GOOS=darwin GOARCH=amd64 $(GO) build -o build/coredns_darwin_amd64
	GOOS=windows GOARCH=amd64 $(GO) build -o build/coredns_windows_amd64.exe
	GOOS=linux GOARCH=amd64 $(GO) build -o build/coredns_linux_amd64
	GOOS=linux GOARCH=arm64 $(GO) build -o build/coredns_linux_arm64
