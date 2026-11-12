# Semana 11 · las seis peticiones del hito 2

| Petición | Con go run | En el contenedor | Con Compose |
| --- | --- | --- | --- |
| 1 sin token | 401 | 401 | |
| 2 login | 200 | 200 | |
| 3 con token | 200 | 200 | |
| 4 rol sin permiso | 403 | 403 | |
| 5 cambio permitido | 200 | 200 | |
| 6 transición prohibida | 422 | 422 | |

Petición 3 con go run (respuesta completa):

```text
listarMascotas OK
```

La misma, en el contenedor (paso 2.5):

```text
listarMascotas OK
```