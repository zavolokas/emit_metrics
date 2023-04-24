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

set-env:
	@set -o allexport; source .env; set +o allexport

get-authz-token: set-env
	@echo 'https://nestservices.google.com/partnerconnections/$(NEST_PROJECT_ID)/auth?redirect_uri=https://www.google.com&access_type=offline&prompt=consent&client_id=$(NEST_CLIENT_ID)&response_type=code&scope=https://www.googleapis.com/auth/sdm.service'

set-tokens:
	# TOKENS_JSON=$(curl -L -X POST "https://www.googleapis.com/oauth2/v4/token?client_id=$(NEST_CLIENT_ID)&client_secret=$(NEST_CLIENT_SECRET)&code=$(AUTHZ_CODE)&grant_type=authorization_code&redirect_uri=https://www.google.com")
	# ACCESS_TOKEN=$(echo "$TOKENS_JSON" | jq -r '.access_token')  
	# REFRESH_TOKEN=$(echo "$TOKENS_JSON" | jq -r '.refresh_token')
	@echo NEST_ACCESS_TOKEN=$(ACCESS_TOKEN) >> .env
	@echo NEST_REFRESH_TOKEN=$(REFRESH_TOKEN) >> .env

