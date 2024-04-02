run:
	@go run main.go
build: 
	@go build -o dist/main main.go
start:
	./dist/main
clean:
		rm -rf dist/
