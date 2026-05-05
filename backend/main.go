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
	db, err = sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	))
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %v", err)
	}
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error al verificar conexión a la base de datos: %v", err)
	}
	log.Println("✓ Conectado a PostgreSQL")
	return nil
}

func setupRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/api/register/cliente", handlers.RegisterClienteHandler).Methods("POST")

	router.HandleFunc("/api/pinturas", handlers.GetPinturasHandler).Methods("GET")
	router.HandleFunc("/api/pinturas/{id}", handlers.GetPinturaByIDHandler).Methods("GET")
	router.HandleFunc("/api/pinturas/artista/{id_artista}", handlers.GetPinturasByArtistaHandler).Methods("GET")
	router.HandleFunc("/api/pinturas/coleccion/{id_coleccion}", handlers.GetPinturasByColeccionHandler).Methods("GET")
	router.HandleFunc("/api/pinturas/tecnica/{id_tecnica}", handlers.GetPinturasByTecnicaHandler).Methods("GET")

	router.HandleFunc("/api/artistas", handlers.GetArtistasHandler).Methods("GET")
	router.HandleFunc("/api/artistas/{id}", handlers.GetArtistaByIDHandler).Methods("GET")

	router.HandleFunc("/api/colecciones", handlers.GetColeccionesHandler).Methods("GET")
	router.HandleFunc("/api/colecciones/{id}", handlers.GetColeccionByIDHandler).Methods("GET")

	router.HandleFunc("/api/tecnicas", handlers.GetTecnicasHandler).Methods("GET")
	router.HandleFunc("/api/tecnicas/{id}", handlers.GetTecnicaByIDHandler).Methods("GET")

	router.HandleFunc("/api/tours", handlers.GetToursHandler).Methods("GET")
	router.HandleFunc("/api/tours/{id}", handlers.GetTourByIDHandler).Methods("GET")

	router.HandleFunc("/api/reportes/pinturas-completo", handlers.ReportePinturasCompletoHandler).Methods("GET")
	router.HandleFunc("/api/reportes/ventas-detalle", handlers.ReporteVentasDetalleHandler).Methods("GET")
	router.HandleFunc("/api/reportes/artistas-resumen", handlers.ReporteArtistasResumenHandler).Methods("GET")

	router.HandleFunc("/api/reportes/artistas-con-ventas", handlers.ReporteArtistasConVentasHandler).Methods("GET")
	router.HandleFunc("/api/reportes/clientes-vip-compradores", handlers.ReporteClientesVIPCompradoresHandler).Methods("GET")

	router.HandleFunc("/api/reportes/ventas-por-mes", handlers.ReporteVentasPorMesHandler).Methods("GET")
	router.HandleFunc("/api/reportes/ventas-por-mes/{anio}", handlers.ReporteVentasPorAnioHandler).Methods("GET")
	router.HandleFunc("/api/reportes/tecnicas-populares", handlers.ReporteTecnicasPopularesHandler).Methods("GET")

	router.HandleFunc("/api/reportes/top-artistas-ventas", handlers.ReporteTopArtistasPorVentasHandler).Methods("GET")
	router.HandleFunc("/api/reportes/colecciones-valor", handlers.ReporteColeccionesValorHandler).Methods("GET")

	router.HandleFunc("/api/exportar/ventas-csv", handlers.ExportarVentasCSVHandler).Methods("GET")
	router.HandleFunc("/api/exportar/pinturas-csv", handlers.ExportarPinturasCSVHandler).Methods("GET")
	router.HandleFunc("/api/exportar/artistas-csv", handlers.ExportarArtistasCSVHandler).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTMiddleware)

	api.HandleFunc("/reservas", handlers.GetReservasHandler).Methods("GET")
	api.HandleFunc("/reservas/{id}", handlers.GetReservaByIDHandler).Methods("GET")
	api.HandleFunc("/reservas", handlers.CreateReservaHandler).Methods("POST")
	api.HandleFunc("/reservas/{id}", handlers.UpdateReservaHandler).Methods("PUT")
	api.HandleFunc("/reservas/{id}", handlers.DeleteReservaHandler).Methods("DELETE")

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

	admin.HandleFunc("/ventas", handlers.GetVentasHandler).Methods("GET")
	admin.HandleFunc("/ventas/{id}", handlers.GetVentaByIDHandler).Methods("GET")
	admin.HandleFunc("/ventas", handlers.CreateVentaHandler).Methods("POST")
	admin.HandleFunc("/ventas/{id}", handlers.UpdateVentaHandler).Methods("PUT")
	admin.HandleFunc("/ventas/{id}", handlers.DeleteVentaHandler).Methods("DELETE")

	admin.HandleFunc("/detalles-venta", handlers.GetDetallesVentaHandler).Methods("GET")
	admin.HandleFunc("/detalles-venta/{id}", handlers.GetDetalleVentaByIDHandler).Methods("GET")
	admin.HandleFunc("/ventas/{id_venta}/detalles", handlers.GetDetallesByVentaHandler).Methods("GET")
	admin.HandleFunc("/detalles-venta", handlers.CreateDetalleVentaHandler).Methods("POST")
	admin.HandleFunc("/detalles-venta/{id}", handlers.UpdateDetalleVentaHandler).Methods("PUT")
	admin.HandleFunc("/detalles-venta/{id}", handlers.DeleteDetalleVentaHandler).Methods("DELETE")

	admin.HandleFunc("/envios", handlers.GetEnviosHandler).Methods("GET")
	admin.HandleFunc("/envios/{id}", handlers.GetEnvioByIDHandler).Methods("GET")
	admin.HandleFunc("/envios", handlers.CreateEnvioHandler).Methods("POST")
	admin.HandleFunc("/envios/{id}", handlers.UpdateEnvioHandler).Methods("PUT")
	admin.HandleFunc("/envios/{id}", handlers.DeleteEnvioHandler).Methods("DELETE")

	admin.HandleFunc("/tours", handlers.CreateTourHandler).Methods("POST")
	admin.HandleFunc("/tours/{id}", handlers.UpdateTourHandler).Methods("PUT")
	admin.HandleFunc("/tours/{id}", handlers.DeleteTourHandler).Methods("DELETE")

	return router
}

func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
