package main

import (
	"encoding/json" // Paquete para manejar el formato JSON
	"fmt"
	"net/http"

	"github.com/Zheta-ai/CloudGuard/internal/database"
)

// Credenciales es el molde (struct) para recibir los datos del usuario desde internet.
// Las etiquetas json:"..." mapean los nombres que vienen de la web a nuestras variables.
type Credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	fmt.Println("🚀 Iniciando los sistemas de CloudGuard...")

	// 1. Conexión a la Bóveda (PostgreSQL)
	db := database.ConectarDB()
	defer db.Close() // Mantiene la puerta abierta hasta que el programa se apague

	// 2. Verificación de Tablas
	database.CrearTablas(db)

	// ================================================================
	// 📡 RUTA API: REGISTRO DE USUARIOS
	// Esta ruta recibirá peticiones externas (JSON) para crear usuarios.
	// ================================================================
	http.HandleFunc("/api/registro", func(w http.ResponseWriter, r *http.Request) {

		// Seguridad básica: Solo permitimos el método POST
		if r.Method != http.MethodPost {
			http.Error(w, "❌ Método no permitido. Usa POST.", http.StatusMethodNotAllowed)
			return
		}

		// Decodificamos el JSON que nos envían
		var creds Credenciales
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "❌ Error al leer los datos JSON", http.StatusBadRequest)
			return
		}

		// Enviamos los datos a la función de registro que ya teníamos
		database.RegistrarUsuario(db, creds.Email, creds.Password)

		// Respondemos al cliente confirmando el éxito
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"mensaje": "✅ Usuario %s registrado con éxito"}`, creds.Email)
	})

	// Ruta de bienvenida/estado de la API
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🌐 CloudGuard API Core - Sistemas Activos")
	})

	// 3. Lanzamiento del Servidor
	fmt.Println("🌐 API escuchando en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("❌ Error crítico al arrancar el servidor:", err)
	}
}
