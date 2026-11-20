//go:build migrar

package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mesa-ayuda/internal/nutricion"
)

func main() {
	databaseURL := os.Getenv("DATABASE_ADMIN_URL")
	if databaseURL == "" {
		log.Fatal("falta la variable de entorno DATABASE_ADMIN_URL")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo conectar con PostgreSQL para migrar: ", err)
	}

	if err := db.AutoMigrate(&nutricion.PlanNutricional{}); err != nil {
		log.Fatal("no se pudo migrar PlanNutricional: ", err)
	}

	log.Println("migración completada correctamente")
}
