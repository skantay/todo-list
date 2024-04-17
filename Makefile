compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans

swag-gen:
	swag init -g ./cmd/app/main.go -o ./docs/api/v1/

test: ### run test
	go clean -testcache
	go test -v ./...

coverage-html: ### run test with coverage and open html report
	go clean -testcache
	go test -coverprofile=cvr.out ./...
	go tool cover -html=cvr.out
	rm cvr.out

coverage: ### run test with coverage
	go clean -testcache
	go test -coverprofile=cvr.out ./...
	go tool cover -func=cvr.out
	rm cvr.out