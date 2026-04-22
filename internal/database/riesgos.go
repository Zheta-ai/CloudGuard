package database

import (
	"database/sql"
	"fmt"
)

// EvaluarTransaccion analiza el riesgo basándose en reglas heurísticas.
func EvaluarTransaccion(db *sql.DB, userID string, monto float64, ubicacion string) (int, string) {
	var score int = 0
	var flags []string

	// =========================================================
	// REGLA 1: Monto Elevado
	// =========================================================
	if monto > 1000.00 {
		score += 50
		flags = append(flags, "high_amount_detected")
	}

	// =========================================================
	// REGLA 2: Velocidad (Velocity Check)
	// Contamos cuántas transacciones existen en el último minuto
	// =========================================================
	var conteoPrevio int
	query := `
		SELECT COUNT(*) 
		FROM transacciones_historial 
		WHERE user_id = $1 
		AND creado_en >= NOW() - INTERVAL '1 minute'
	`

	err := db.QueryRow(query, userID).Scan(&conteoPrevio)
	if err != nil {
		fmt.Println("⚠️ Error al consultar ráfaga:", err)
	}

	// Si ya existen 3 o más transacciones en el último minuto, esta dispara el riesgo.
	if conteoPrevio >= 3 {
		score += 40
		flags = append(flags, "high_velocity_detected")
	}

	// =========================================================
	// REGISTRO: Guardamos la transacción actual
	// PostgreSQL asignará la hora exacta (UTC) automáticamente
	// =========================================================
	insertQuery := `
		INSERT INTO transacciones_historial (user_id, monto, ubicacion) 
		VALUES ($1, $2, $3)
	`
	_, err = db.Exec(insertQuery, userID, monto, ubicacion)
	if err != nil {
		fmt.Println("⚠️ Error al guardar en historial:", err)
	}

	// Limitar el score máximo a 100
	if score > 100 {
		score = 100
	}

	// Formatear flags para la respuesta
	resumenFlags := "ok"
	if len(flags) > 0 {
		resumenFlags = ""
		for i, f := range flags {
			if i > 0 {
				resumenFlags += " | "
			}
			resumenFlags += f
		}
	}

	return score, resumenFlags
}
