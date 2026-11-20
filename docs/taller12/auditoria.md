# Semana 12 · Auditoría de inyección

Búsqueda ejecutada: `git grep -n -E -f docs/taller12/patrones.txt -- "*.go" ":(exclude)*_test.go"` (salida en `docs/taller12/auditoria.txt`, 10 líneas).

| Línea | ¿Arma SQL? | ¿Llega texto del cliente? | ¿Va aparte? | Qué hicimos |
| --- | --- | --- | --- | --- |
| `internal/almacen/filas/consultas.sql.go:19` (`row := q.db.QueryRowContext(ctx, buscarPlan, id)`) | Sí | Sí (`id`) | Sí (argumento separado) | Revisada: consulta parametrizada, no concatena SQL. |
| `internal/almacen/filas/consultas.sql.go:52` (`result, err := q.db.ExecContext(ctx, guardarPlan, ...`) | Sí | Sí (campos del plan) | Sí (parámetros posicionales) | Revisada: escritura parametrizada con `ExecContext`. |
| `internal/nutricion/repositorio.go:72` (`err := r.DB.First(&plan, id).Error`) | Sí | Sí (`id`) | Sí (consulta ORM por clave primaria) | Revisada: consulta parametrizada con GORM `First`. |
| `internal/almacen/filas/db.go:13, 15, 16` (`ExecContext`, `QueryContext`, `QueryRowContext`) | No (firmas de interfaz) | N/A | Sí (`...interface{}`) | Revisada: interfaces estándar de consultas con argumentos variádicos. |
| `internal/almacen/almacen.go:74`, `internal/nutricion/manejadores.go:37`, `internal/tickets/tickets.go:115`, `internal/tickets/tickets.go:195` | No (ruta `.tmp`, grupo de rutas chi, lectura de query `estado` y mensaje de error HTTP en `respuesta.Error`) | Parcial (`estado`, `t.Estado`, `body.Estado` en respuesta/filtro) | No aplica SQL (`tickets.go:195` concatena con `+` el texto del error HTTP `transicion_invalida`, no una consulta SQL) | Revisada: se verificó que ninguna de estas líneas arma SQL y que `tickets.go:195` solo construye el mensaje de error HTTP. |

**Resultado:** Todas las líneas de `auditoria.txt` fueron revisadas. Las consultas SQL/ORM pasan parámetros por separado y la concatenación con `+` en `internal/tickets/tickets.go:195` corresponde al mensaje de error HTTP (`respuesta.Error`), sin inyección SQL.