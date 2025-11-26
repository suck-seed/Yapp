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



.PHONY: buildBackend startBackend startNoLogsBackend removeContainer removeContainerData startApp
