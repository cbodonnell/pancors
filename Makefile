build:
	go build -o ./bin/pancors cmd/pancors/main.go

serve:
	./bin/pancors

container:
	docker build -t cheebz/pancors .
