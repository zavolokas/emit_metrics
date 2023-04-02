run:
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go

docker-build:
	docker build . -t metric-emitter:v1.0.0

docker-run:
	docker run -d --env-file .env --restart=unless-stopped  --network=host --name metric-emitter metric-emitter:v1.0.0

docker-logs:
	docker logs --follow metric-emitter
