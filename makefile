binaryName?=ipam
version?=0.6

build:
	GOOS=darwin GOARCH=arm64 go build -o ./releases/arm64/mac/$(binaryName) -ldflags="-X 'main.version=v$(version)-DARWIN-ARM64'" *.go
	# GOOS=darwin GOARCH=amd64 go build -o ./releases/x86_64/mac/$(binaryName) -ldflags="-X 'main.version=v$(version)-DARWIN-AMD64'" *.go
	# GOOS=linux GOARCH=amd64 go build -o ./releases/x86_64/linux/$(binaryName) -ldflags="-X 'main.version=v$(version)-AMD64'" *.go
	# GOOS=windows GOARCH=amd64 go build -o ./releases/x86_64/windows/$(binaryName).exe -ldflags="-X 'main.version=v$(version)-Microsoft-Windows'" *.go
