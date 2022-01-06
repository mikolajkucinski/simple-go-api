.SILENT:

compile:
	protoc --go_out=./internal/proto-files/ ./internal/proto-files/users.proto
	go build -o . ./cmd/main.go

run:
	echo "Running program..."
	./main