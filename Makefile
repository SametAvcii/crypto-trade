build:
	go build -o crypto-trade  main.go

run: 
	./crypto-trade

swagger:
	swag init --quiet --parseDependency

hot:
	docker compose -p crypto-trade -f docker-compose-hot.yml up

run:
	docker compose -f docker-compose.yml up -d --build

down:
	docker compose -f docker-compose.yml down