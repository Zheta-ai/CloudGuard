package database

import (
	"database/sql"
	"fmt"
	"log"

	// Importamos nuestra nueva trituradora
	"golang.org/x/crypto/bcrypt"
)

// RegistrarUsuario toma el email y la contraseña limpia, la cifra, y la guarda en PostgreSQL.
func RegistrarUsuario(db *sql.DB, email string, passwordLimpia string) {
	// 1. Cifrar (Hashear) la contraseña
	// Transformamos el texto plano en un código indescifrable
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordLimpia), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Error crítico al cifrar la contraseña: ", err)
	}

	// 2. Preparar la inyección a la base de datos
	// NOTA PRO DE CIBERSEGURIDAD: Usamos $1 y $2 en lugar de pegar el texto directamente.
	// Esto previene los ataques de "SQL Injection" (que hackers borren tu base de datos).
	query := `INSERT INTO usuarios (email, password_hash) VALUES ($1, $2)`

	// 3. Ejecutar la orden
	// string(hash) convierte los bytes triturados en texto para guardarlos
	_, err = db.Exec(query, email, string(hash))
	if err != nil {
		// Si el usuario ya existe (porque el email es UNIQUE), saldrá este mensaje sin tumbar el servidor
		fmt.Println("⚠️  Atención al registrar (¿el usuario ya existe?):", err)
		return
	}

	fmt.Println("👤 ¡Nuevo usuario registrado y blindado con éxito! ->", email)
}

// VerificarUsuario busca a un usuario por su email y extrae sus datos de la tabla.
func VerificarUsuario(db *sql.DB, email string) {
	// 1. Preparamos la búsqueda (SELECT)
	// Le pedimos a PostgreSQL que nos traiga estas 3 columnas específicas
	query := `SELECT id, email, password_hash FROM usuarios WHERE email = $1`

	// 2. Variables vacías donde Go guardará los resultados
	var id int
	var emailDB string
	var hashDB string

	// 3. Ejecutamos la búsqueda
	// QueryRow busca una sola fila. Scan() agarra los datos de la base de datos
	// y los "inyecta" en nuestras variables vacías usando los punteros (&).
	err := db.QueryRow(query, email).Scan(&id, &emailDB, &hashDB)
	if err != nil {
		fmt.Println("❌ No se encontró al usuario:", err)
		return
	}

	// 4. Imprimimos el reporte clasificado
	fmt.Println("--------------------------------------------------")
	fmt.Println("🔍 REPORTE DE SEGURIDAD (BÓVEDA DE POSTGRESQL)")
	fmt.Printf("ID: %d\n", id)
	fmt.Printf("Email: %s\n", emailDB)
	fmt.Printf("Contraseña original: [DESTRUIDA/NO VISIBLE]\n")
	fmt.Printf("Lo que realmente se guardó (Hash): %s\n", hashDB)
	fmt.Println("--------------------------------------------------")
}
