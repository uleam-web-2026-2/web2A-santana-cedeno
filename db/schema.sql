CREATE TABLE plan_nutricionals (
    id          bigserial PRIMARY KEY,
    id_mascota  bigint NOT NULL,
    estado      text NOT NULL,
    descripcion text NOT NULL,
    fecha_inicio text NOT NULL,
    fecha_fin    text NOT NULL
);