# Hito 2 - Hoja de trabajo

## 1 · Usuarios de prueba

| Usuario | Rol         | Contraseña de ejemplo |
|---------|-------------|-----------------------|
| ana     | cliente     | ana-2026              |
| luis    | veterinario | luis-2026             |


---

## 2 · Las seis peticiones

| # | Petición                          | Con qué credencial          | Debe responder |
|---|-----------------------------------|-----------------------------|----------------|
| 1 | GET /mascotas                     | Ninguna                     | 401            |
| 2 | POST /login con ana               | Ninguna                     | 200 y el token |
| 3 | GET /mascotas                     | Token de Ana (cliente)      | 200            |
| 4 | PATCH /planes/{id} → activo       | Token de Ana (cliente)      | 403            |
| 5 | PATCH /planes/{id} → activo       | Token de Luis (veterinario) | 200            |
| 6 | PATCH /planes/{id} → finalizado   | Token de Luis (veterinario) | 422            |

---

## 3 · Matriz ruta × rol

| Ruta              | Sin sesión | Cliente | Veterinario |
|-------------------|------------|---------|-------------|
| POST /login       | Sí         | Sí      | Sí          |
| GET /mascotas     | No (401)   | Sí      | Sí          |
| POST /mascotas    | No (401)   | Sí      | Sí          |
| GET /planes/{id}  | No (401)   | Sí      | Sí          |
| PATCH /planes/{id}| No (401)   | No (403)| Sí          |
| POST /detalles    | No (401)   | (plan)  | (plan)      |
| POST /citas       | No (401)   | (plan)  | (plan)      |

---

## 4 · Recuadro de cumplimiento

☑ Los structs de todas las entidades están en el código y sus tablas existen  
☑ El validador termina sin errores ni avisos  
☑ Las seis peticiones están en la colección, en orden  
☑ Cada una responde lo que dice esta hoja  
☑ `go test ./...` en verde, con una prueba de integración  
☑ El código de la entidad principal está en capas y la estructura está explicada en `docs/decisiones.md`  
☑ Hay una captura de respaldo de cada paso
