compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
