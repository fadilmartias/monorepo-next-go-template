SERVICE=backend

up:
	docker compose up -d

up-build:
	docker compose up -d --build

build:
	docker compose build

down:
	docker compose down

down-v:
	docker compose down -v

restart:
	docker compose restart

ps:
	docker compose ps

logs:
	docker compose logs -f --tail=500

logs:
	docker compose logs -f --tail=500 ${SERVICE}

exec-root:
	docker compose exec -u root ${SERVICE} sh

exec:
	docker compose exec ${SERVICE} sh

pull:
	docker compose pull

clean:
	docker system prune -f