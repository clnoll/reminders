WITH_ENV=env $$(xargs < env.sh)
SHELL = bash -u

run-server:
	$(WITH_ENV) go run api/main.go

run-worker:
	$(WITH_ENV) go run worker/main.go

test:
	go test ./...
