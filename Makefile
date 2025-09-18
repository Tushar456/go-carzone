.PHONY: up down logs restart


up:
	docker-compose up --build -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart


