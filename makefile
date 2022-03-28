run:
	go run main.go

buildDocker:
	docker build -t my-go-server .