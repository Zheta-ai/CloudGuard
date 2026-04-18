package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConectarDB() *sql.DB {
	connStr := "user=admin password=supersecretpassword dbname=cloudguard_db port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Error al configurar la base de datos: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("❌ No se pudo conectar a PostgreSQL: ", err)
	}
	fmt.Println("✅ ¡Base de datos conectada con éxito, papu!")
	return db
}

func CrearTablas(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("❌ Error crítico al crear las tablas en la DB: ", err)
	}
	fmt.Println("🛡️  Estructura de la base de datos: VERIFICADA Y LISTA.")
}
