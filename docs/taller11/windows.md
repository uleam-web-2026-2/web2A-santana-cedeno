# Semana 11 · Incidencias en Windows

## 1. Preparación de WSL 2

### Qué falló

Al comenzar el taller, el equipo no tenía preparado WSL 2 y Docker Desktop no podía utilizarse correctamente para ejecutar los contenedores Linux requeridos por el taller.

### Por qué ocurrió

Docker Desktop en Windows puede utilizar WSL 2 como backend para ejecutar contenedores Linux. En nuestro equipo fue necesario habilitar las características correspondientes de Windows antes de continuar.

### Qué hicimos

Se habilitaron las características necesarias de WSL y de la plataforma de máquina virtual desde PowerShell con permisos de administrador.

Después de reiniciar Windows se comprobó:

- WSL instalado.
- Versión predeterminada de WSL configurada como 2.
- Virtualización habilitada.
- Docker Desktop funcionando.

Finalmente se ejecutó:

`docker run --rm hello-world`

y se obtuvo el mensaje `Hello from Docker!`, confirmando que Docker funcionaba correctamente.

## 2. Redirección de salida en PowerShell

### Qué falló

Al guardar algunas evidencias de comandos de Docker utilizando `> archivo 2>&1` desde Windows PowerShell, el archivo generado incluía información adicional con el texto `NativeCommandError`.

Esto ocurrió, por ejemplo, al guardar la salida del contenedor sin variables y los logs del contenedor.

### Por qué ocurrió

La salida real del programa se estaba escribiendo correctamente, pero Windows PowerShell añadía información adicional al manejar la salida de error de algunos comandos externos.

### Qué hicimos

Para generar las evidencias limpias ejecutamos la redirección mediante `cmd /c`.

Por ejemplo:

`cmd /c "docker run --rm nutricion > docs\taller11\1_sin_variables.txt 2>&1"`

y:

`cmd /c "docker logs nutricion-api > docs\taller11\contenedor_log.txt 2>&1"`

De esta manera los archivos conservaron la salida real del contenedor sin el bloque adicional de `NativeCommandError`.

## 3. Conexión de PostgreSQL desde Docker

### Qué falló

Al ejecutar la API dentro del contenedor con `host=localhost`, la conexión a PostgreSQL falló con `connection refused`.

### Por qué ocurrió

Dentro del contenedor, `localhost` representa al propio contenedor y no al equipo Windows donde estaba ejecutándose PostgreSQL.

### Qué hicimos

En el archivo local `.env.docker` se cambió:

`host=localhost`

por:

`host=host.docker.internal`

Se mantuvo el puerto `5432`, que es el puerto utilizado por PostgreSQL en nuestro proyecto.

Después de este cambio, el contenedor pudo conectarse correctamente a la base `nutricion_mascotas` y atender las seis peticiones del Hito 2.