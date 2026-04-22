package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// ConectarDB se queda EXACTAMENTE IGUAL como la tienes.
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

// CAMBIO AQUÍ: Modificamos la tabla que se crea
func CrearTablas(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS transacciones_historial (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(50) NOT NULL,
		monto DECIMAL(10, 2) NOT NULL,
		ubicacion VARCHAR(100),
		creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("❌ Error crítico al crear las tablas en la DB: ", err)
	}
	fmt.Println("🛡️  Estructura de Riesgo (Transacciones): VERIFICADA Y LISTA.")
}
