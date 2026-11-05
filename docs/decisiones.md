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