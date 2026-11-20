# Semana 12 · CORS

## Configuración aplicada

Se configuró CORS como middleware superior, inmediatamente después de crear el router con `chi.NewRouter()`.

Para las pruebas del taller se definió como único origen permitido:

`http://localhost:5173`

No se utilizó `*`.

## Origen permitido

Se realizó una petición previa `OPTIONS` desde:

`http://localhost:5173`

para comprobar el método `PATCH`.

La respuesta fue `200 OK` e incluyó:

`Access-Control-Allow-Origin: http://localhost:5173`

También indicó:

`Access-Control-Allow-Methods: PATCH`

La salida completa está guardada en:

`docs/taller12/cors_permitido.txt`

## Origen no permitido

Se repitió la misma petición previa utilizando:

`http://ejemplo-no-permitido.test`

La respuesta no incluyó la cabecera:

`Access-Control-Allow-Origin`

Por tanto, ese origen no recibió permiso CORS.

La salida completa está guardada en:

`docs/taller12/cors_no_permitido.txt`

## Petición curl sin Origin

También se realizó una petición normal con `curl` a:

`GET /tickets`

sin enviar la cabecera `Origin`.

La API respondió:

`200 OK`

Esto demuestra que CORS controla qué orígenes puede utilizar un navegador para acceder a las respuestas, pero una petición normal realizada sin `Origin` continúa llegando a la API.

La salida completa está guardada en:

`docs/taller12/cors_sin_origin.txt`

## Resultado

Las pruebas realizadas muestran que:

- `http://localhost:5173` recibe permiso CORS.
- Un origen diferente no recibe `Access-Control-Allow-Origin`.
- Una petición normal de `curl` sin `Origin` sigue funcionando.