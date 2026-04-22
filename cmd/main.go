package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Zheta-ai/CloudGuard/internal/database"
)

// 1. NUEVOS STRUCTS: Lo que recibimos y lo que respondemos
type PeticionTransaccion struct {
	UserID    string  `json:"user_id"`
	Monto     float64 `json:"monto"`
	Ubicacion string  `json:"ubicacion"`
}

type RespuestaRiesgo struct {
	RiskScore int    `json:"risk_score"`
	Status    string `json:"status"` // allow, verify, block
	Flag      string `json:"flag"`
}

func main() {
	fmt.Println("🚀 Iniciando CloudGuard Risk-API...")

	// 1. Conexión a la Bóveda y creación de la tabla de transacciones
	db := database.ConectarDB()
	defer db.Close()
	database.CrearTablas(db)

	// ================================================================
	// 📡 RUTA API: ANÁLISIS DE RIESGO EN TIEMPO REAL
	// ================================================================
	http.HandleFunc("/api/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "❌ Método no permitido. Usa POST.", http.StatusMethodNotAllowed)
			return
		}

		// Decodificamos el JSON que nos envía la Fintech
		var peticion PeticionTransaccion
		err := json.NewDecoder(r.Body).Decode(&peticion)
		if err != nil {
			http.Error(w, "❌ Error al leer los datos JSON", http.StatusBadRequest)
			return
		}

		// MANDAMOS A EVALUAR AL MOTOR DE RIESGO (nuestra función en riesgo.go)
		puntaje, flagRecibido := database.EvaluarTransaccion(db, peticion.UserID, peticion.Monto, peticion.Ubicacion)

		// Decidimos la acción (Status) basados en el puntaje
		var statusFinal string
		if puntaje >= 50 {
			statusFinal = "block"
		} else if puntaje >= 20 {
			statusFinal = "verify"
		} else {
			statusFinal = "allow"
		}

		// Armamos la respuesta
		respuesta := RespuestaRiesgo{
			RiskScore: puntaje,
			Status:    statusFinal,
			Flag:      flagRecibido,
		}

		// Enviamos el JSON de vuelta al cliente
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(respuesta)
	})

	// Ruta de estado
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🌐 CloudGuard Risk-API - Sistemas Activos")
	})

	// Lanzamiento del Servidor
	fmt.Println("🌐 API escuchando en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("❌ Error crítico al arrancar el servidor:", err)
	}
}
