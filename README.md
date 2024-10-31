# golang-backend

## Docker set up

Necesitamos crear una red externa para posteriormente hacer llamadas a la API desde la aplicacion de astro

    docker network create suffgo-network

Levantamos el contenedor

    docker compose up 

## Ejecutar migraciones

    sh /scripts/migrate.sh
    
## Endpoints

El formato estandar de los endpoints va a ser `v1/{entidad}` (solo para endpoints que tengan que ver con los datos)

Create user POST `localhost:{API_PORT}/v1/users`

    body 
    {
        "name": "Tiago",
        "lastname": "Cardenas",
        "username": "tiaguinho",
        "dni": "14123122",
        "email": "tcardenas@gmail.com",
        "password": "gaturro01"
    },
    {
        "name": "Ignacio",
        "lastname": "Sanchez",
        "username": "neich",
        "dni": "43299985",
        "email": "nachoagusss1@@gmail.com",
        "password": "goesgodcolapinto"
    },
    {
        "name": "Marcos",
        "lastname": "Soto",
        "username": "msoto",
        "dni": "33998222",
        "email": "m.soto123@@gmail.com",
        "password": "redhot"
    },
    {
        "name": "Constanza",
        "lastname": "Benedetti",
        "username": "coty2",
        "dni": "12312312",
        "email": "coreanos@@gmail.com",
        "password": "cotoyoteconozco"
    },
     {
        "name": "Cristian",
        "lastname": "Balihaut",
        "username": "crisbal",
        "dni": "339988798",
        "email": "b.cris@@gmail.com",
        "password": "racingmicorazon"
    }


Retrieve all users GET `localhost:{API_PORT}/v1/users`

Retrieve user by id GET `localhost:{API_PORT}/v1/users/:id`

Delete user by id DELETE `localhost:{API_PORT}/v1/users/:id`

