# Decisiones de diseño

## 1. Orden de los middleware
- **Qué se decidió**: Se colocó primero el middleware de recuperación de pánico y luego el de registro.
- **Por qué**: Así se asegura que cualquier error inesperado se capture y se devuelva en JSON, y además quede registrado en el log.
- **Alternativa descartada**: Poner el logger primero y el recovery después, lo cual hacía que el log no se escribiera en caso de pánico.
- **Cómo se verificó**: Se probó con `curl -i localhost:8080/explotar` y se confirmó que el servidor devolvió 500 en JSON y el log anotó la petición.

---

## 2. Distribución de códigos de error
- **Qué se decidió**: Cada situación devuelve un código específico:  
  - 201 → creación de tickets válidos  
  - 200 → consultas exitosas  
  - 404 → recurso o ruta inexistente  
  - 405 → método no permitido  
  - 400 → cuerpo malformado o id inválido  
  - 422 → datos inválidos  
  - 500 → pánico controlado
- **Por qué**: Mantener consistencia con la guía y facilitar la depuración.
- **Alternativa descartada**: Usar siempre 400 para cualquier error.
- **Cómo se verificó**: Se ejecutaron los `curl` de Fase 1 y Fase 2, confirmando que cada caso devuelve el código correcto y siempre en JSON.

---

## 3. Manejo de campos desconocidos en el body
- **Qué se decidió**: Si el cliente envía un campo extra en el JSON, el servidor devuelve error `datos_invalidos`.
- **Por qué**: Mantener un contrato estricto de la API y evitar datos inesperados que puedan romper la lógica.
- **Alternativa descartada**: Ignorar los campos desconocidos y procesar solo los válidos.
- **Cómo se verificó**: Se probó enviando un JSON con un campo `"extra":"valor"` y se recibió el error esperado en formato JSON con código 422.

---

## Criterio del hito
Ni una sola respuesta en texto plano. Todas las respuestas son JSON con el código correcto.

## HITO 2


## 4. Organización en capas de la entidad principal
- **Qué se decidió**: Separar responsabilidades en archivos distintos:
  - `manejadores.go`: define las rutas y recibe las peticiones HTTP.
  - `estado.go`: contiene la lógica de negocio para validar transiciones de estado.
  - `modelos.go`: define las entidades con GORM y JSON.
  - `usuarios.go`, `tokens.go`, `guardia.go`: manejan login, emisión de tokens y middleware de seguridad.
  - `respuesta.go`: unifica el formato de salida (Error, Exito).
- **Por qué**: Evita mezclar responsabilidades, facilita pruebas unitarias y cumple con el principio de capas.
- **Alternativa descartada**: Poner toda la lógica en un solo archivo.
- **Cómo se verificó**: Se ejecutaron las pruebas (`go test ./...`) y se comprobó que las rutas responden en JSON con el código correcto.


# Decisiones de arquitectura – Hito 2

## 1. Organización por capas

Para el desarrollo del Hito 2 se organizó el módulo de nutrición en tres responsabilidades principales:

- **Manejador HTTP:** recibe las solicitudes, valida los datos de entrada y devuelve las respuestas HTTP correspondientes.
- **Servicio:** contiene las reglas de negocio, especialmente la validación de las transiciones de estado de los planes nutricionales.
- **Repositorio:** se encarga de consultar y guardar los planes nutricionales. Se implementaron repositorios en memoria y mediante GORM.

## 2. Separación de responsabilidades

El manejador HTTP no realiza consultas directas a la base de datos. En su lugar, utiliza el servicio de planes, que trabaja mediante la interfaz `RepositorioPlanes`.

La regla de negocio para cambiar el estado es independiente del protocolo HTTP. Esto permite reutilizarla y probarla sin depender de una solicitud web.

## 3. Validación de transiciones

Se establecieron las siguientes transiciones permitidas:

| Estado actual | Estados permitidos |
|---|---|
| Pendiente | Activo, cancelado |
| Activo | Finalizado, cancelado |
| Finalizado | Ninguno |
| Cancelado | Ninguno |

Cuando se intenta realizar una transición no permitida, el sistema devuelve el código HTTP 422.

## 4. Autenticación y autorización

Se implementó autenticación mediante tokens JWT. Las rutas protegidas requieren un token válido.

También se incorporó control de acceso basado en roles para restringir determinadas operaciones a usuarios autorizados.

## 5. Pruebas realizadas

Se ejecutó `go test ./...` para comprobar los paquetes del proyecto.

Además, se verificaron las siguientes respuestas:

- HTTP 401: solicitud a una ruta protegida sin autenticación válida.
- HTTP 403: usuario sin el rol requerido.
- HTTP 422: transición de estado no permitida.
- HTTP 200: transición de estado permitida.

Se agregó una prueba de integración para comprobar que una transición no permitida devuelve HTTP 422 utilizando autenticación JWT.

## 6. Consideraciones

Para la demostración actual se utiliza un repositorio en memoria porque el servidor se ejecuta con `DB: nil`. También se implementó un repositorio GORM para permitir el uso de una base de datos cuando exista una conexión configurada.


---

## Semana 10 · transacciones

### D1 · El caso

Actualmente nuestro proyecto de nutrición para mascotas no tiene un caso de uso implementado que necesite realizar dos escrituras en dos tablas o dos filas como una única operación.

Por este motivo, siguiendo las indicaciones del taller, no se añadieron tablas ni operaciones artificiales solamente para implementar una transacción.

El problema planteado en el taller con ConfirmarYCobrar sirve para mostrar qué ocurre cuando una operación necesita realizar dos escrituras relacionadas. Si la primera escritura se guarda correctamente y la segunda falla, el sistema queda en un estado inconsistente aunque el cliente reciba un error.

En nuestro proyecto podría aparecer una situación equivalente si en el futuro se implementara una operación para crear un plan nutricional junto con su detalle nutricional. La operación tendría que guardar el PlanNutricional y el DetalleNutricional como una sola unidad de trabajo.

Si se guardara el plan y después fallara la creación del detalle, quedaría un plan nutricional incompleto. Por tanto, ambas escrituras deberían completarse o deshacerse juntas.

### D2 · Opciones que consideramos

Consideramos dos alternativas.

**Opción 1: controlar cada escritura por separado desde el servicio.**

El servicio podría llamar primero al repositorio para guardar el plan y después realizar una segunda llamada para guardar su detalle.

Esta alternativa mantiene una implementación sencilla, pero no garantiza atomicidad. Si la primera operación termina correctamente y la segunda falla, la primera modificación ya habría quedado almacenada.

**Opción 2: delegar la transacción al adaptador del repositorio.**

El servicio expresaría la operación de negocio mediante la interfaz del repositorio, sin conocer GORM, database/sql ni los detalles de la transacción.

Cada adaptador sería responsable de ejecutar las escrituras dentro de una transacción de la tecnología que utiliza. De esta manera, si una de las operaciones falla, el adaptador podría deshacer el conjunto completo.

Esta opción respeta la separación entre la lógica de negocio y la infraestructura.

### D3 · La que elegimos y por qué

Elegiríamos que la transacción sea responsabilidad del adaptador del repositorio.

La razón en nuestro proyecto es que ServicioPlanes debe mantenerse independiente de la tecnología utilizada para persistir los datos. Actualmente ya trabajamos mediante la interfaz RepositorioPlanes y durante el taller incorporamos un adaptador que utiliza sqlc para determinadas operaciones.

Por ello, el servicio no debería conocer cómo comienza, confirma o revierte una transacción.

Si en el futuro la creación de un plan nutricional y sus detalles tuviera que ser atómica, la interfaz representaría esa operación de negocio y cada adaptador resolvería internamente su propia transacción.

### D4 · La que descartamos y por qué

Descartamos gestionar las dos escrituras de manera independiente desde el servicio.

En nuestro caso, el problema sería que podría guardarse un PlanNutricional y fallar posteriormente la creación de su DetalleNutricional. El sistema devolvería un error al cliente, pero el plan ya habría quedado almacenado y la operación estaría completada solo parcialmente.

También descartamos que el servicio gestione directamente una transacción de GORM o database/sql, porque eso haría que la lógica de negocio dependiera de una tecnología concreta de persistencia.

### D5 · Con sqlc

En la Fase 3 incorporamos sqlc mediante un adaptador que implementa la misma interfaz utilizada por ServicioPlanes.

Para un futuro caso transaccional, mantendríamos el mismo principio: el servicio no tendría que modificarse para conocer sqlc. La transacción sería resuelta dentro del adaptador correspondiente.

En el taller actual no añadimos consultas ni transacciones nuevas para este punto porque nuestro proyecto todavía no posee un caso de uso implementado que realice las dos escrituras relacionadas. Implementarlas solamente para esta fase supondría inventar una operación que actualmente no existe en el proyecto.

### D6 · Archivos

Para esta fase se modificó:

- `docs/decisiones.md`

No fue necesario modificar:

- `internal/nutricion/modelos.go`
- `internal/nutricion/servicio.go`
- `internal/nutricion/repositorio.go`
- `internal/nutricion/repositorio_sqlc.go`
- `internal/nutricion/manejadores.go`
- `db/schema.sql`
- `db/consultas.sql`
- `main.go`

La razón es que esta fase se resolvió como una decisión de diseño, ya que el proyecto no dispone actualmente de un caso de uso que necesite dos escrituras atómicas.

---

## Semana 11 · Docker

### D1 · Imagen de dos etapas

Para contenerizar la API se utilizó un Dockerfile de dos etapas.

La primera etapa utiliza `golang:1.27` y contiene las herramientas necesarias para descargar las dependencias y compilar el proyecto. En esta etapa se ejecuta la compilación con `CGO_ENABLED=0` para obtener el binario `/servidor`.

La segunda etapa utiliza `gcr.io/distroless/static-debian13:nonroot` y recibe únicamente el binario generado en la etapa anterior. De esta forma, la imagen final no necesita incluir todo el entorno de compilación de Go.

La diferencia se comprobó en `docs/taller11/imagenes.txt`: nuestra imagen `nutricion:compilacion` ocupa aproximadamente 1.62 GB, mientras que la imagen final `nutricion:latest` ocupa aproximadamente 38.7 MB.

También se comprobó en `docs/taller11/capas.txt` que el binario `/servidor` ocupa aproximadamente 21.3 MB dentro de la imagen.

El archivo `.dockerignore` impide que archivos sensibles como `.env` y la carpeta `.git` entren en el contexto de construcción. Esto se comprobó mediante `docs/taller11/contexto.txt`.

Finalmente, `docs/taller11/usuario.txt` devuelve `65532`, comprobando que la imagen final ejecuta el proceso con el usuario sin privilegios `nonroot`.

### D2 · Variables y secretos

La configuración que cambia entre máquinas no se incluyó dentro de la imagen Docker. Nuestra API necesita principalmente `DATABASE_URL` para conectarse a PostgreSQL y `TOKEN_SECRETO` para trabajar con los tokens.

Se comprobó este comportamiento ejecutando la imagen sin proporcionar variables de entorno. El contenedor se detuvo indicando que faltaba `TOKEN_SECRETO`, y el resultado quedó guardado en `docs/taller11/1_sin_variables.txt`.

Para ejecutar el contenedor utilizamos un archivo local `.env.docker`. Este archivo contiene los valores necesarios para nuestra máquina, pero no se sube al repositorio porque `.gitignore` excluye `.env.*`.

El archivo `.env.example` sí permanece en el repositorio porque únicamente muestra la estructura de configuración y utiliza valores de ejemplo, sin contener secretos reales.

De esta manera, la misma imagen Docker puede recibir la configuración necesaria al momento de arrancar sin guardar contraseñas ni secretos directamente dentro de la imagen.

### D3 · Conexión desde el contenedor

Al ejecutar inicialmente el contenedor con `DATABASE_URL` configurada con `host=localhost`, la conexión con PostgreSQL falló.

La evidencia quedó guardada en `docs/taller11/2_local.txt`, donde se muestra que el contenedor intentó conectarse a `127.0.0.1:5432` y recibió `connection refused`.

El problema ocurre porque el contenedor tiene su propio entorno de red. Por tanto, `localhost` dentro del contenedor no representa la máquina Windows donde se está ejecutando PostgreSQL.

Para resolverlo en Docker Desktop cambiamos únicamente el host de la conexión en `.env.docker`:

`host=localhost`

por:

`host=host.docker.internal`

Se mantuvieron los demás datos de conexión de nuestro proyecto, incluida la base `nutricion_mascotas` y el puerto `5432`.

Después del cambio, el contenedor pudo conectarse correctamente a PostgreSQL y la API quedó ejecutándose y accesible mediante el puerto `8080`.

### D4 · Comportamiento de la API dentro del contenedor

Después de proporcionar las variables de entorno correctas, configurar `host.docker.internal` para acceder a PostgreSQL y publicar el puerto `8080`, la API pudo ejecutarse completamente dentro del contenedor Docker.

Se repitieron las mismas seis peticiones utilizadas como línea base con `go run` y se obtuvieron los mismos códigos HTTP:

| Petición | Con go run | En el contenedor |
| --- | --- | --- |
| Sin token | 401 | 401 |
| Login | 200 | 200 |
| Con token | 200 | 200 |
| Rol sin permiso | 403 | 403 |
| Cambio permitido | 200 | 200 |
| Transición prohibida | 422 | 422 |

También se comparó la respuesta completa de la petición 3. Tanto con `go run` como desde el contenedor se obtuvo:

`listarMascotas OK`

Los resultados quedaron registrados en `docs/taller11/peticiones.md` y el funcionamiento del contenedor también quedó reflejado en `docs/taller11/contenedor_log.txt`.

Por tanto, contenerizar la API no cambió el comportamiento observable de las seis peticiones utilizadas como referencia.

### D5 · API y PostgreSQL con Docker Compose

Para el reto se creó `compose.yaml` con dos servicios: `api` y `db`.

El servicio `api` construye la imagen utilizando el `Dockerfile` del proyecto y publica el puerto `8080` para que la API siga siendo accesible desde `http://127.0.0.1:8080`.

El servicio `db` utiliza PostgreSQL y crea la base `nutricion_mascotas`.

Con Docker Compose ya no fue necesario utilizar `host.docker.internal`. La variable `DATABASE_URL` de la API utiliza `host=db`, porque `db` es el nombre del servicio de PostgreSQL dentro de la configuración de Compose.

De esta forma, un solo comando permite iniciar la API y la base de datos como servicios de la misma aplicación.

### D6 · Comunicación entre servicios

Con Docker Compose, la API y PostgreSQL forman parte de la misma aplicación y pueden comunicarse utilizando el nombre del servicio definido en `compose.yaml`.

Por esta razón, la conexión de la API utiliza:

`host=db`

en lugar de:

`host=localhost`

o:

`host=host.docker.internal`

El nombre `db` corresponde al servicio PostgreSQL definido en `compose.yaml`.

Esta configuración permite que la API se conecte directamente al servicio de base de datos administrado por Docker Compose, sin depender del PostgreSQL instalado directamente en Windows.

Se comprobó el funcionamiento ejecutando `docker compose up -d`. PostgreSQL quedó en estado `healthy` y la API arrancó correctamente. Los registros de ambos servicios quedaron guardados en `docs/taller11/compose_log.txt`.

### D7 · Persistencia de PostgreSQL

En `compose.yaml` se definió el volumen llamado `nutricion_datos` para almacenar los datos de PostgreSQL.

El servicio `db` monta este volumen en:

`/var/lib/postgresql/data`

La finalidad es separar los datos de la base de datos del ciclo de vida del contenedor. De esta manera, los datos no dependen únicamente del sistema de archivos interno del contenedor de PostgreSQL.

Docker administra los volúmenes como almacenamiento persistente y Docker Compose puede reutilizar un volumen existente al volver a crear los servicios.

Para nuestro proyecto, esto permite conservar la información de `nutricion_mascotas` aunque se recree el contenedor de PostgreSQL.

### D8 · Orden de arranque y estado saludable

PostgreSQL necesita estar preparado para aceptar conexiones antes de que la API intente utilizar la base de datos.

Por este motivo, en `compose.yaml` se añadió un `healthcheck` al servicio `db` utilizando `pg_isready`.

El servicio `api` utiliza:

`depends_on`

con la condición:

`service_healthy`

De esta manera, Docker Compose espera a que PostgreSQL alcance el estado saludable antes de iniciar la API.

Durante la prueba se comprobó que el servicio `db` apareció como `healthy` y posteriormente el servicio `api` quedó iniciado correctamente.

Los registros almacenados en `docs/taller11/compose_log.txt` también muestran que PostgreSQL quedó preparado para aceptar conexiones y que la API pudo atender las peticiones posteriormente.

### D9 · Comprobación final con Docker Compose

La solución del reto se comprobó levantando toda la aplicación mediante:

`docker compose up -d`

Con este comando se iniciaron los dos servicios definidos en `compose.yaml`: la API y PostgreSQL.

La API siguió disponible en `http://127.0.0.1:8080` y se repitieron las seis peticiones utilizadas durante el taller.

Los resultados obtenidos con Docker Compose fueron:

| Petición | Con Compose |
| --- | --- |
| Sin token | 401 |
| Login | 200 |
| Con token | 200 |
| Rol sin permiso | 403 |
| Cambio permitido | 200 |
| Transición prohibida | 422 |

Estos códigos coinciden con los obtenidos anteriormente mediante `go run` y mediante el contenedor individual.

Los resultados quedaron registrados en `docs/taller11/peticiones.md` y los registros de la API y PostgreSQL quedaron guardados en `docs/taller11/compose_log.txt`.

Por tanto, el comportamiento observable de la API se mantuvo al ejecutar la API y PostgreSQL juntos mediante Docker Compose.

---

## Semana 12 · Seguridad
### D1 · Auditoría de inyección

Se realizó una búsqueda del código del proyecto utilizando los patrones definidos en `docs/taller12/patrones.txt`.

La salida se guardó en `docs/taller12/auditoria.txt` y se revisaron las coincidencias encontradas para distinguir cuáles estaban relacionadas con consultas SQL y cuáles correspondían únicamente a operaciones normales de Go.

Las consultas generadas por sqlc pasan los valores como argumentos separados y el repositorio GORM utiliza un identificador de tipo `uint` en `First`, sin construir la consulta mediante concatenación de texto.

También se comprobó que el filtro `estado` encontrado en el módulo de tickets trabaja sobre un repositorio en memoria y no construye una consulta SQL.

Como comprobación práctica se envió a la ruta real `/tickets/{id}` un valor equivalente a `' OR 1=1 --`. La API respondió `400 Bad Request` con el código `id_invalido`, sin devolver tickets ni información adicional.

La respuesta real de esta prueba quedó guardada en `docs/taller12/inyeccion.txt`.

### D2 · Política CORS

Se configuró CORS como middleware superior del router, inmediatamente después de `chi.NewRouter()`.

Para las pruebas del taller se definió únicamente el origen:

`http://localhost:5173`

No se utilizó `*`.

Se realizaron tres comprobaciones:

1. Una petición previa `OPTIONS` desde `http://localhost:5173` respondió `200 OK` e incluyó `Access-Control-Allow-Origin: http://localhost:5173`.

2. La misma petición desde `http://ejemplo-no-permitido.test` no incluyó la cabecera `Access-Control-Allow-Origin`.

3. Una petición normal realizada con `curl` sin cabecera `Origin` continuó funcionando y respondió `200 OK`.

Las salidas reales quedaron guardadas en:

- `docs/taller12/cors_permitido.txt`
- `docs/taller12/cors_no_permitido.txt`
- `docs/taller12/cors_sin_origin.txt`

De esta manera, la API no utiliza un origen abierto con `*` y solo concede autorización CORS al origen configurado para las pruebas del taller.

### D3 · Límite de intentos en el login

Se protegió únicamente la ruta pública `POST /login` con un límite de 5 solicitudes por minuto por IP.

La IP del cliente se obtiene desde la conexión mediante `ClientIPFromRemoteAddr`. Después, la función `porIP` utiliza esa IP como clave del limitador.

Para comprobar la defensa se realizaron seis intentos consecutivos de inicio de sesión utilizando una contraseña incorrecta.

Los primeros cinco intentos llegaron al login y respondieron `401 Unauthorized`.

El sexto intento fue detenido por el limitador y respondió:

`429 Too Many Requests`

Además, la respuesta incluyó:

`Retry-After: 60`

La salida real de los seis intentos quedó guardada en:

`docs/taller12/limite_login.txt`

Después de esta prueba se reinició el servidor y se volvieron a ejecutar las seis peticiones del Hito 2 para comprobar que el límite del login no alteró el comportamiento esperado del resto de la API.

### D4 · Comprobación de las defensas

Después de aplicar las defensas de inyección, CORS y límite de intentos en el login, se volvieron a ejecutar las seis peticiones utilizadas como línea base del Hito 2.

Los resultados fueron:

| Petición | Antes (taller 11) | Con las defensas |
| --- | --- | --- |
| Sin token | 401 | 401 |
| Login | 200 | 200 |
| Con token | 200 | 200 |
| Rol sin permiso | 403 | 403 |
| Cambio permitido | 200 | 200 |
| Transición prohibida | 422 | 422 |

Los mismos códigos HTTP se mantuvieron después de aplicar las nuevas defensas.

Los resultados quedaron registrados en `docs/taller12/peticiones.md`.

Además, el límite del login se comprobó independientemente mediante seis intentos consecutivos con credenciales incorrectas. Los primeros cinco llegaron al login y el sexto fue detenido con `429 Too Many Requests` y la cabecera `Retry-After`.

De esta manera, las nuevas defensas añadidas en la Semana 12 no modificaron el comportamiento esperado de las seis peticiones utilizadas como referencia.

---

## Semana 12 · reto: mínimo privilegio en la base

### D5 · Con qué usuario se conecta la API

Para la ejecución normal de la API se creó el usuario de PostgreSQL `nutricion_api`.

Este usuario no es superusuario, no puede crear bases de datos ni roles y no es propietario de la tabla `plan_nutricionals`. La tabla continúa perteneciendo al usuario administrador `postgres`.

Se concedieron únicamente los permisos necesarios para las operaciones actuales de la API:

- `CONNECT` sobre la base `nutricion_mascotas`.
- `USAGE` sobre el esquema `public`.
- `SELECT`, `INSERT` y `UPDATE` sobre `plan_nutricionals`.
- `USAGE` y `SELECT` sobre la secuencia `plan_nutricionals_id_seq`.

No se concedieron permisos `DELETE`, `CREATE` ni `DROP`.

Para comprobar el mínimo privilegio se intentó borrar `plan_nutricionals` utilizando `nutricion_api` dentro de una transacción. PostgreSQL respondió:

`ERROR: must be owner of table plan_nutricionals`

La evidencia real quedó guardada en `docs/taller12/reto_drop.txt`.

También se guardó `docs/taller12/reto_usuarios.txt`, donde se observa que `postgres` conserva los privilegios administrativos y `nutricion_api` no tiene atributos especiales.

Finalmente, se arrancó la API utilizando `nutricion_api` y las seis peticiones conservaron los códigos esperados: `401`, `200`, `200`, `403`, `200` y `422`.

### D6 · Quién crea y cambia las tablas

La ejecución normal de la API ya no realiza `AutoMigrate`.

Se separó la responsabilidad de modificar el esquema de la base de datos mediante `migrar.go`. Este archivo utiliza la variable `DATABASE_ADMIN_URL` y se ejecuta con:

`go run -tags migrar .`

Para esta operación se utiliza el usuario administrador `postgres`, que continúa siendo propietario de `plan_nutricionals`.

La API normal utiliza `DATABASE_URL` y se conecta como `nutricion_api`, que no es propietario de la tabla y no tiene permisos para crear, borrar ni modificar su estructura.

Se comprobó la migración por separado y el comando respondió:

`migración completada correctamente`

Por tanto, si en el futuro cambia el modelo `PlanNutricional`, primero debe ejecutarse la migración con el usuario administrador. Arrancar únicamente la API con `nutricion_api` no debe modificar automáticamente el esquema de PostgreSQL.

### D7 · Lo que esto no frena

El principio de mínimo privilegio reduce el daño que podría causar un problema de seguridad en la API, pero no sustituye las demás defensas.

Si en el futuro existiera una inyección SQL que consiguiera ejecutar operaciones con el usuario `nutricion_api`, ese usuario seguiría pudiendo realizar las operaciones para las que tiene permiso:

- consultar registros de `plan_nutricionals`;
- insertar nuevos planes nutricionales;
- actualizar planes existentes;
- utilizar la secuencia necesaria para generar identificadores.

Sin embargo, `nutricion_api` no tiene permisos administrativos ni es propietario de la tabla.

La prueba realizada demostró que no puede ejecutar:

`DROP TABLE plan_nutricionals`

PostgreSQL respondió:

`ERROR: must be owner of table plan_nutricionals`

Por tanto, el mínimo privilegio limita el alcance de un posible fallo, pero no evita por sí mismo una inyección SQL. Por esa razón se mantienen también las consultas parametrizadas y las demás defensas implementadas durante el taller.