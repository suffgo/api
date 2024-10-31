# golang-backend

## Docker set up

Necesitamos crear una red externa para posteriormente hacer llamadas a la API desde la aplicacion de astro

    docker network create suffgo-network

Levantamos el contenedor

    docker compose up 

## Ejecutar migraciones

docker exec -it go-app bash

go run cmd/migrate/main.go

## Endpoints

El formato estandar de los endpoints va a ser `v1/{entidad}` (solo para endpoints que tengan que ver con los datos)

Create user POST `localhost:{API_PORT}/v1/users`

    body 
    {
        "name":"Tiago",
        "lastname":"Cardenas",
        "username": "tiaguinho",
        "dni": "1412312",
        "email": "tcardenas@gmail.com",
        "Password": "gaturro01"
    }


Retrieve all users GET `localhost:{API_PORT}/v1/users`

Retrieve user by id GET `localhost:{API_PORT}/v1/users/:id`

Delete user by id DELETE `localhost:{API_PORT}/v1/users/:id`

