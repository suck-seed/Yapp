buildBackend:
	docker compose up --build

startBackend:
	docker compose up

startNoLogsBackend:
	docker compose up -d


removeContainer:
	docker compose down


removeContainerData:
	docker compose down --volumes






.PHONY: build start startNoLogs removeContainer removeContainerData
