# Guia rapida para iniciar API para desarrollo.

## Docker set up

Necesitamos crear una red externa para posteriormente hacer llamadas a la API desde la aplicacion de astro.

    docker network create suffgo-network
### Enviroment

    cp .env.sample .env

Completar las variables de entorno necesarias.

### Sesiones

**importante!!** Definir llave con la que se firman las cookies de sesion en .env.

    SECRET_SESSION_AUTH_KEY=

Crear llave secreta.

    openssl rand -hex 32

> Necesitan tener openssl instalado, por defecto ya se encuentra instalado en linux.


Levantamos el contenedor.

    docker compose up 

Ejecutar migraciones

    sh scripts/migrate.sh

> En produccion son automaticas.

## Troubleshooting

Si desea restaurar la base de datos puede borrar la carpeta docker que se encuentra en raiz.
