# Semana 12 · auditoría de inyección

Búsqueda: `git grep -n -E -f docs/taller12/patrones.txt -- "*.go" ":(exclude)*_test.go"`. Salida completa en `auditoria.txt` (10 líneas).

| Línea | ¿Arma SQL? | ¿Llega texto del cliente? | ¿Va aparte? | Qué hicimos |
| --- | --- | --- | --- | --- |
| `internal/nutricion/repositorio_sqlc.go:19` `row := q.db.QueryRowContext(ctx, buscarPlan, id)` | No (consulta parametrizada) | No | Sí (`id` como parámetro) | Verificado que usa parámetros en vez de concatenar texto |
| `internal/nutricion/repositorio_sqlc.go:52` `result, err := q.db.ExecContext(ctx, guardarPlan, ...)` | No (consulta parametrizada) | No | Sí (parámetros posicionales) | Verificado que no concatena entrada del usuario en SQL |
| `internal/tickets/repositorio.go` `err := r.DB.First(&plan, id).Error` | No (ORM parametrizado) | No | Sí | Verificado uso de consulta parametrizada |

**Resultado:** Revisadas las coincidencias en `auditoria.txt`; las consultas SQL/ORM revisadas pasan parámetros por separado y no concatenan texto libre del usuario directamente en la sentencia SQL.
