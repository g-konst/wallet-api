.DEFAULT_GOAL := up

up:
	docker-compose up --build

test:
	go test ./tests
