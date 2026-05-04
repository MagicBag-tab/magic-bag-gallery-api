package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"magic-bag-gallery-api/internal/handlers"
	"magic-bag-gallery-api/internal/middleware"

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

	router.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/api/register/cliente", handlers.RegisterClienteHandler).Methods("POST")

	router.HandleFunc("/api/pinturas", handlers.GetPinturasHandler).Methods("GET")
	router.HandleFunc("/api/pinturas/{id}", handlers.GetPinturaByIDHandler).Methods("GET")
	router.HandleFunc("/api/artistas", handlers.GetArtistasHandler).Methods("GET")
	router.HandleFunc("/api/artistas/{id}", handlers.GetArtistaByIDHandler).Methods("GET")
	router.HandleFunc("/api/colecciones", handlers.GetColeccionesHandler).Methods("GET")
	router.HandleFunc("/api/colecciones/{id}", handlers.GetColeccionByIDHandler).Methods("GET")
	router.HandleFunc("/api/tecnicas", handlers.GetTecnicasHandler).Methods("GET")
	router.HandleFunc("/api/tecnicas/{id}", handlers.GetTecnicaByIDHandler).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)

	admin := router.PathPrefix("/api").Subrouter()
	admin.Use(middleware.JWTMiddleware)
	admin.Use(middleware.RequireRole("empleado"))

	admin.HandleFunc("/auth/register/empleado", handlers.RegisterEmpleadoHandler).Methods("POST")

	admin.HandleFunc("/pinturas", handlers.CreatePinturaHandler).Methods("POST")
	admin.HandleFunc("/pinturas/{id}", handlers.UpdatePinturaHandler).Methods("PUT")
	admin.HandleFunc("/pinturas/{id}", handlers.DeletePinturaHandler).Methods("DELETE")

	admin.HandleFunc("/artistas", handlers.CreateArtistaHandler).Methods("POST")
	admin.HandleFunc("/artistas/{id}", handlers.UpdateArtistaHandler).Methods("PUT")
	admin.HandleFunc("/artistas/{id}", handlers.DeleteArtistaHandler).Methods("DELETE")

	admin.HandleFunc("/colecciones", handlers.CreateColeccionHandler).Methods("POST")
	admin.HandleFunc("/colecciones/{id}", handlers.UpdateColeccionHandler).Methods("PUT")
	admin.HandleFunc("/colecciones/{id}", handlers.DeleteColeccionHandler).Methods("DELETE")

	admin.HandleFunc("/tecnicas", handlers.CreateTecnicaHandler).Methods("POST")
	admin.HandleFunc("/tecnicas/{id}", handlers.UpdateTecnicaHandler).Methods("PUT")
	admin.HandleFunc("/tecnicas/{id}", handlers.DeleteTecnicaHandler).Methods("DELETE")

	admin.HandleFunc("/usuarios", handlers.GetUsuariosHandler).Methods("GET")
	admin.HandleFunc("/usuarios/{id}", handlers.GetUsuarioByIDHandler).Methods("GET")
	admin.HandleFunc("/usuarios", handlers.CreateUsuarioHandler).Methods("POST")
	admin.HandleFunc("/usuarios/{id}", handlers.UpdateUsuarioHandler).Methods("PUT")
	admin.HandleFunc("/usuarios/{id}", handlers.DeleteUsuarioHandler).Methods("DELETE")

	//router.HandleFunc("/api/pinturaByArtista/{id_artista}", handlers.GetPinturasByArtistaHandler).Methods("GET")
	//router.HandleFunc("/api/pinturaByColeccion/{id_coleccion}", handlers.GetPinturasByColeccionHandler).Methods("GET")
	//router.HandleFunc("/api/pinturaByTecnica/{id_tecnica}", handlers.GetPinturasByTecnicaHandler).Methods("GET")

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
