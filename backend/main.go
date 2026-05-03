package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"magic-bag-gallery-api/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	loadDatabase()
	handlers.SetDB(db)

	router := setupRouter()

	handler := corsMiddleware()(router)

	log.Println("Servidor iniciado en http://localhost:8888")
	if err := http.ListenAndServe(":8888", handler); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func loadDatabase() {
	if err := loadEnv(); err != nil {
		log.Fatalf("Error al cargar variables de entorno: %v", err)
	}

	if err := connectDB(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
}

func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error al cargar variables de entorno: %v", err)
	}
	return nil
}

func connectDB() error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB")))
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %v", err)
	}
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error al verificar conexión a la base de datos: %v", err)
	}
	return nil
}

func setupRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/pinturas", handlers.GetPinturasHandler).Methods("GET")
	router.HandleFunc("/api/pinturas/{id}", handlers.GetPinturaByIDHandler).Methods("GET")
	router.HandleFunc("/api/pinturas", handlers.CreatePinturaHandler).Methods("POST")
	router.HandleFunc("/api/pinturas/{id}", handlers.UpdatePinturaHandler).Methods("PUT")
	router.HandleFunc("/api/pinturas/{id}", handlers.DeletePinturaHandler).Methods("DELETE")

	//router.HandleFunc("/api/artistas", getArtistas).Methods("GET")
	//router.HandleFunc("/api/artistas/{id}", getArtistaByID).Methods("GET")
	//router.HandleFunc("/api/artistas", createArtista).Methods("POST")
	//router.HandleFunc("/api/artistas/{id}", updateArtista).Methods("PUT")
	//router.HandleFunc("/api/artistas/{id}", deleteArtista).Methods("DELETE")

	//router.HandleFunc("/api/colecciones", getColecciones).Methods("GET")
	//router.HandleFunc("/api/colecciones/{id}", getColeccionByID).Methods("GET")
	//router.HandleFunc("/api/colecciones", createColeccion).Methods("POST")
	//router.HandleFunc("/api/colecciones/{id}", updateColeccion).Methods("PUT")
	//router.HandleFunc("/api/colecciones/{id}", deleteColeccion).Methods("DELETE")

	//router.HandleFunc("/api/tecnicas", getTecnicas).Methods("GET")
	//router.HandleFunc("/api/tecnicas/{id}", getTecnicaByID).Methods("GET")
	//router.HandleFunc("/api/tecnicas", createTecnica).Methods("POST")
	//router.HandleFunc("/api/tecnicas/{id}", updateTecnica).Methods("PUT")
	//router.HandleFunc("/api/tecnicas/{id}", deleteTecnica).Methods("DELETE")

	//router.HandleFunc("/api/usuarios", getUsuarios).Methods("GET")
	//router.HandleFunc("/api/usuarios/{id}", getUsuarioByID).Methods("GET")
	//router.HandleFunc("/api/usuarios", createUsuario).Methods("POST")
	//router.HandleFunc("/api/usuarios/{id}", updateUsuario).Methods("PUT")
	//router.HandleFunc("/api/usuarios/{id}", deleteUsuario).Methods("DELETE")

	//router.HandleFunc("/api/ventas", getVentas).Methods("GET")
	//router.HandleFunc("/api/ventas/{id}", getVentaByID).Methods("GET")
	//router.HandleFunc("/api/ventas", createVenta).Methods("POST")
	//router.HandleFunc("/api/ventas/{id}", updateVenta).Methods("PUT")
	//router.HandleFunc("/api/ventas/{id}", deleteVenta).Methods("DELETE")

	//router.HandleFunc("/api/envios", getEnvios).Methods("GET")
	//router.HandleFunc("/api/envios/{id}", getEnvioByID).Methods("GET")
	//router.HandleFunc("/api/envios", createEnvio).Methods("POST")
	//router.HandleFunc("/api/envios/{id}", updateEnvio).Methods("PUT")
	//router.HandleFunc("/api/envios/{id}", deleteEnvio).Methods("DELETE")

	//router.HandleFunc("/api/tours", getTours).Methods("GET")
	//router.HandleFunc("/api/tours/{id}", getTourByID).Methods("GET")
	//router.HandleFunc("/api/tours", createTour).Methods("POST")
	//router.HandleFunc("/api/tours/{id}", updateTour).Methods("PUT")
	//router.HandleFunc("/api/tours/{id}", deleteTour).Methods("DELETE")

	//router.HandleFunc("/api/reservas", getReservas).Methods("GET")
	//router.HandleFunc("/api/reservas/{id}", getReservaByID).Methods("GET")
	//router.HandleFunc("/api/reservas", createReserva).Methods("POST")
	//router.HandleFunc("/api/reservas/{id}", updateReserva).Methods("PUT")
	//router.HandleFunc("/api/reservas/{id}", deleteReserva).Methods("DELETE")

	//router.HandleFunc("/api/pinturaByTecnica/{id_tecnica}", getPinturasByTecnicaHandler).Methods("GET")
	//router.HandleFunc("/api/pinturaByColeccion/{id_coleccion}", getPinturasByColeccionHandler).Methods("GET")
	//router.HandleFunc("/api/pinturaByArtista/{id_artista}", getPinturasByArtistaHandler).Methods("GET")

	return router
}

func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
