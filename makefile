run:
	@go run main.go

build: 
	@go build -o dist/main main.go

start:
	./dist/main

migration_up:
	migrate -path ./database/migrations -database "${DATABASE_URL}" up

migration_down:
	migrate -path ./database/migrations -database "${DATABASE_URL}" down

clean:
		rm -rf dist/
