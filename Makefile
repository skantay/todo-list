compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans

linter-golangci: ### check by golangci linter
	golangci-lint run

test: ### run test
	go test -v ./...

coverage-html: ### run test with coverage and open html report
	go test -coverprofile=cvr.out ./...
	go tool cover -html=cvr.out
	rm cvr.out

coverage: ### run test with coverage
	go test -coverprofile=cvr.out ./...
	go tool cover -func=cvr.out
	rm cvr.out

swag-gen:
	swag init -g ./cmd/main.go -o ./docs/api/v1/