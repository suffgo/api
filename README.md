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


### Endpoints protegidos

Login POST `localhost:{API_PORT}/v1/users/login`

body 
    {
        "username":"yourusername",
        "password":"yourpassword"
    }

> En caso exitoso, devuelve un token

Supongamos que ahora necesito acceder a un endpoint que esta protegido con autenticacion JWT como por ejemplo GET `localhost:{API_PORT}/secure`

En el header tenes que crear un campo Authorization con el siguiente valor:

    Bearer {token-que-te-devuelve-la-ruta-login}

Y te va a devolver un mensaje

{  
  "message": "Hola max, usted esta autorizado"
}