
DOCKER_COMPOSE = docker compose

COMPOSE_FILE = compose.yml
ENV_FILE = .env.local
DOCKER_COMPOSE_ENV_LOCAL = $(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(ENV_FILE)

up: docker-up
up-d: docker-up-d
down: docker-down
restart: down up
ps: docker-ps
logs: docker-logs

docker-up:
	$(DOCKER_COMPOSE_ENV_LOCAL) up --build --remove-orphans

docker-up-d:
	$(DOCKER_COMPOSE_ENV_LOCAL) up -d --build --remove-orphans

docker-down:
	$(DOCKER_COMPOSE_ENV_LOCAL) down --remove-orphans

docker-down-clear:
	$(DOCKER_COMPOSE_ENV_LOCAL) down -v --remove-orphans

docker-ps:
	$(DOCKER_COMPOSE_ENV_LOCAL) ps --all

docker-logs:
	$(DOCKER_COMPOSE_ENV_LOCAL) logs

.PHONY: tests