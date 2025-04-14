build:
	go build -o crypto-trade  main.go

swagger:
	swag init --generalInfo cmd/app/main.go

hot:
	docker compose -p crypto-trade -f docker-compose-hot.yml up

run:
	docker compose -f docker-compose.yml up -d --build

down:
	docker compose -f docker-compose.yml down