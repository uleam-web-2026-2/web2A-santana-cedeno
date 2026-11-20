#!/usr/bin/env bash
set -e

psql -v ON_ERROR_STOP=1 \
  --username "$POSTGRES_USER" \
  --dbname "$POSTGRES_DB" <<-EOSQL

CREATE TABLE IF NOT EXISTS plan_nutricionals (
    id BIGSERIAL PRIMARY KEY,
    id_mascota BIGINT NOT NULL,
    estado TEXT NOT NULL,
    descripcion TEXT NOT NULL,
    fecha_inicio TEXT NOT NULL,
    fecha_fin TEXT NOT NULL
);

CREATE ROLE nutricion_api
    WITH LOGIN
    PASSWORD '${API_DB_PASSWORD}';

GRANT CONNECT
    ON DATABASE nutricion_mascotas
    TO nutricion_api;

GRANT USAGE
    ON SCHEMA public
    TO nutricion_api;

GRANT SELECT, INSERT, UPDATE
    ON TABLE plan_nutricionals
    TO nutricion_api;

GRANT USAGE, SELECT
    ON SEQUENCE plan_nutricionals_id_seq
    TO nutricion_api;

EOSQL
