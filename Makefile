build: 
	go build ./cmd/main.go

run: 
	go run ./cmd/main.go conf/example.yml

test:
	go test -v ./...