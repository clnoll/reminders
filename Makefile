WITH_ENV=env $$(xargs < env.sh)
SHELL = bash -u

run-server: export ENV = DEV
run-server:
	$(WITH_ENV) go run api/main.go

run-worker: export ENV = DEV
run-worker:
	$(WITH_ENV) go run worker/main.go

temporal:
	docker-compose -f docker-compose/docker-compose.yml up


test: export ENV = TEST
test:
	go test ./...
