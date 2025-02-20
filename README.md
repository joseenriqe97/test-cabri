# Cabri Test

> Este proyecto es una API RESTful construida con Golang/echo que permite la creacion y autentificacion de usuario.

## Comenzando 🚀

Estas instrucciones te permitirán tener una copia del proyecto ejecutándose en tu máquina local para desarrollo y pruebas.

### Pre-requisitos 📋

Antes de empezar, asegúrate de tener instalado lo siguiente:

* **Docker**:  [Descarga e Instalación de Docker](https://docs.docker.com/get-docker/)
* **Docker Compose**:  Generalmente se instala con Docker Desktop. Si necesitas instalarlo por separado, consulta la [Documentación de Docker Compose](https://docs.docker.com/compose/install/)

### Instalación y Ejecución 🔧

1. **Clona el repositorio**
   ```bash
   git clone [https://github.com/joseenriqe97/test-cabri.git)
   cd [nombre del repositorio]

2. **Exportar las variables en tu terminal (para sesiones temporales)**

    ```bash
        export DATABASE_USER=postgres
        export DATABASE_PASSWORD=0000
        export DATABASE_DB_NAME=cabri
        export AUTH_TOKEN=assswww
        export URL_PLD=[http://98.81.235.22](http://98.81.235.22)

3. **Ejecuta la aplicación con Docker Compose**

    ```bash
        docker-compose up --build
        
        OR

        DATABASE_USER=postgres DATABASE_PASSWORD=0000 DATABASE_DB_NAME=cabri AUTH_TOKEN=AUTH URL_PLD=http://98.81.235.22 docker-compose up --build

