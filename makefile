.PHONY: dev prod down build clean

dev:
	docker compose --profile dev up --build

prod:
	docker compose --profile prod up --build

down:
	docker compose down

build-dev:
	docker compose --profile dev --build

build-prod:
	docker compose --profile prod --build

clean:
	docker compose down --volumes --rmi all