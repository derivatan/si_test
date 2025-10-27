## Run integration tests
integration:
	docker compose -f docker-compose.yml up --abort-on-container-exit --exit-code-from runner

## Cleanup
clean:
	docker compose -f docker-compose.yml rm

db:
	docker compose -f docker-compose.yml start database

