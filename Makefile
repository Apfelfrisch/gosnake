build-native:
	@echo "Building..."
	@go build -o ./build/snake ./runner/ebitengine/
build-windows:
	@echo "Building for Windows..."
	@env GOOS=windows GOARCH=amd64 go build -o ./build/snake.exe ./runner/ebitengine/
build-mac-arm:
	@echo "Building for Mac Arms..."
	@env GOOS=darwin GOARCH=arm64 go build -o ./build/snake-arm ./runner/ebitengine/
run: build-native
	./build/snake
test:
	@go test ./...
