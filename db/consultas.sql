-- name: BuscarPlan :one
SELECT id, id_mascota, estado, descripcion, fecha_inicio, fecha_fin
FROM plan_nutricionals
WHERE id = $1;

-- name: GuardarPlan :execrows
UPDATE plan_nutricionals
SET id_mascota = $2,
    estado = $3,
    descripcion = $4,
    fecha_inicio = $5,
    fecha_fin = $6
WHERE id = $1;