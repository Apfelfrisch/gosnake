build-native:
	@echo "Building..."
	@go build -o ./build/snake ./runner/ebitengine/
build-windows:
	@echo "Building for Windows..."
	@env GOOS=windows GOARCH=amd64 go build -o ./build/snake.exe ./runner/ebitengine/
run: build-native
	./build/snake
test:
	@go test ./...
