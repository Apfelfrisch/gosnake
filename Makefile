build-native:
	@echo "Building..."
	@go build -o ./build/snake ./cmd/
build-windows:
	@echo "Building for Windows..."
	@env GOOS=windows GOARCH=amd64 go build -o ./build/snake.exe ./cmd/
build-mac-arm:
	@echo "Building for Mac Arms..."
	@env GOOS=darwin GOARCH=arm64 go build -o ./build/snake-arm ./cmd/
run: build-native
	./build/snake
proto:
	@protoc ./game/network/payload/payload.proto --go_out=.
test:
	@go test ./...
