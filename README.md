# golang-backend

## Docker set up

Necesitamos crear una red externa para posteriormente hacer llamadas a la API desde la aplicacion de astro

    docker network create suffgo-network

Levantamos el contenedor

    docker compose up 

## Ejecutar migraciones

docker exec -it go-app bash

go run cmd/migrate/main.go