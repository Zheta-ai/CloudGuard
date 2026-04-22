package database

import (
	"database/sql"
	"fmt"
)

// EvaluarTransaccion revisa los datos y devuelve un puntaje de riesgo (0 a 100)
// y un mensaje de bandera (flag).
func EvaluarTransaccion(db *sql.DB, userID string, monto float64, ubicacion string) (int, string) {
	var score int = 0
	var flag string = "ok"

	// REGLA HEURÍSTICA 1: Si el monto es anormalmente alto para una transacción estándar
	if monto > 1000.00 {
		score += 50
		flag = "high_amount_detected"
	}

	// 2. Guardamos la transacción en el historial para futuras evaluaciones de ráfaga
	query := `INSERT INTO transacciones_historial (user_id, monto, ubicacion) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, userID, monto, ubicacion)
	if err != nil {
		fmt.Println("⚠️ Error guardando el historial:", err)
	}

	return score, flag
}
