buildBackend:
	docker compose up --build

startBackend:
	docker compose up

startNoLogsBackend:
	docker compose up -d


stopBackend:
	docker compose down


removeContainerData:
	docker compose down --volumes


startApp:
	docker compose up -d
	cd yapp-frontend && npm run dev

docs:
	swag init -g cmd/yapppp-server/main.go --output docs

run: docs
	go run cmd/yapppp-server/main.go


.PHONY: buildBackend startBackend startNoLogsBackend removeContainerData startApp docs run
