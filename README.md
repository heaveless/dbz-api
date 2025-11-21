# DBZ API (Go + Gin)

![Hero](https://github.com/heaveless/dbz-api/blob/master/assets/hero.png?raw=true)

Este proyecto expone una **API REST** para obtener información de personajes (por ejemplo, personajes de Dragon Ball u otro universo ficticio).

La aplicación está construida con:

- **Golang** (Go) como lenguaje principal.
- **Gin** como framework HTTP.
- **MongoDB** como base de datos principal.
- Un **cliente HTTP** con *circuit breaker* y *fallback* hacia una **API externa** de personajes.

Además, implementa una lógica de **fallback**:

1. Primero intenta obtener el personaje desde la **base de datos**.
2. Si no existe (o hay un error que cumple ciertas condiciones), consulta una **API externa** usando un cliente HTTP protegido con **circuit breaker**.
3. Si la API externa responde correctamente, el personaje se **guarda en la base de datos** y se devuelve al cliente.
4. Si la llamada externa también falla, se devuelve un **error controlado** al cliente.

## Índice

- [DBZ API (Go + Gin)](#dbz-api-go--gin)
- [1. Requisitos previos](#1-requisitos-previos)
- [2. Clonar el proyecto](#2-clonar-el-proyecto)
- [3. Configuración de entorno](#3-configuración-de-entorno)
  - [3.1. Variables mínimas necesarias](#31-variables-mínimas-necesarias)
- [4. Levantar infraestructura (MongoDB con Docker)](#4-levantar-infraestructura-mongodb-con-docker)
  - [4.1. Archivo docker-compose.yml de ejemplo](#41-archivo-docker-composeyml-de-ejemplo)
  - [4.2. Levantar Contenedores](#42-levantar-contenedores)
- [5. Endpoints de la API](#5-endpoints-de-la-api)
  - [5.1. Obtener personaje por nombre](#51-obtener-personaje-por-nombre)
  - [8.2. Respuesta esperada](#82-respuesta-esperada)
- [9. Arquitectura en capas](#9-arquitectura-en-capas)
- [10. Diagrama de secuencias](#10-diagrama-de-secuencias)

---

## 1. Requisitos previos

Antes de ejecutar el proyecto, asegúrate de tener instalado:

- [Go](https://go.dev/dl/) **1.25+**
- [Git](https://git-scm.com/)
- (Opcional) [Docker](https://www.docker.com/) y [Docker Compose](https://docs.docker.com/compose/) para levantar MongoDB fácilmente.
- (Opcional) `make` si quieres usar comandos predefinidos en un `Makefile`.

---

## 2. Clonar el proyecto

```bash
git clone https://github.com/heaveless/dbz-api.git
```

## 3. Configuración de entorno

La aplicación usa variables de entorno para configurar su comportamiento. Puedes:

- Crear un archivo .env en la raíz del proyecto

- Exportar las variables directamente en tu terminal/shell.

### 3.1. Variables mínimas necesarias

Ejemplo de archivo .env:

```bash
APP_ENV=development
APP_PORT=4000

DB_HOST=mongodb
DB_PORT=27017
DB_NAME=mydb

API_URI=https://dragonball-api.com
```

## 4. Levantar infraestructura (MongoDB con Docker)

Si no tienes MongoDB instalado localmente, puedes usar Docker para levantarlo rápidamente.

### 4.1. Archivo docker-compose.yml de ejemplo

Crea un archivo docker-compose.yml en la raíz del proyecto:

```yml
version: "3.8"

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    image: app
    container_name: app
    restart: unless-stopped
    env_file: .env
    ports:
      - "$APP_PORT:$APP_PORT"
    depends_on:
      - mongodb

  mongodb:
    image: mongo:6.0
    container_name: mongodb
    restart: unless-stopped
    env_file: .env
    ports:
      - "$DB_PORT:$DB_PORT"
    volumes:
      - mongodb_data:/data/db

volumes:
  mongodb_data:
```

### 4.2. Levantar Contenedores

```bash
docker-compose up -d
```
Esto creará:

- Un contenedor llamado `mongodb`.
- Un contenedor llamado `app`
- App escuchando en el puerto `4000`.
- Un volumen `mongo_data` para persistir los datos.

Asegúrate de que `APP_PORT` en tu `.env` coincide con esta configuración.

La API estará disponible en:

```bash
http://localhost:4000
```

## 5. Endpoints de la API

Actualmente, la API expone al menos un endpoint principal para consultar personajes.

### 5.1. Obtener personaje por nombre

**Objetivo**: Devolver la información de un personaje a partir de su nombre.

- **Método**: GET
- **Path**: characters
- **Body**:
    - `name` (string, **requerido**): nombre del personaje que se desea consultar.

**Request de ejemplo (cURL)**:

```bash
curl -X POST "http://localhost:4000/characters" \
  -H "Content-Type: application/json" \
  -d '{ "name": "Goku" }'
```

### 8.2. Respuesta esperada

```json
{
  "data": {
    "id": 1,
    "name": "Goku",
    "race": "Saiyan",
    ...
  }
}
```

## 9. Arquitectura en capas
Estamos aplicando una arquitectura en capas inspirada en **Clean Architecture** / **Hexagonal**, separando claramente:

- **Capa de entrada (Delivery)** → HTTP / transport
- **Capa de Aplicación** → casos de uso / servicios
- **Capa de Dominio** → entidades, interfaces y reglas de negocio
- **Capa de Infraestructura** → implementaciones técnicas (DB, HTTP, circuit breakers)
- **Bootstrap** → composición de dependencias
- **Utilidades transversales** → helpers genéricos (fallback, etc.)

![Clean Architecture Diagram](https://github.com/heaveless/dbz-api/blob/master/assets/architecture.jpg?raw=true)

## 10. Diagrama de secuencias

![Secuence Diagram](https://github.com/heaveless/dbz-api/blob/master/assets/secuence.png?raw=true)