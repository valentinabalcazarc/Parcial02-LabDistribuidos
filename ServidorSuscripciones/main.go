package main

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"pasarela-pagos/config"
	"pasarela-pagos/controllers"
	"pasarela-pagos/data/repositories"
	"pasarela-pagos/middleware"
	"pasarela-pagos/services"
)

func main() {
	cfg := config.Cargar()

	// Capa de acceso a datos
	suscripcionRepo := repositories.NuevoSuscripcionRepository()
	tokenRepo, err := repositories.NuevoTokenRepository(filepath.Join(cfg.DirectorioDatos, "tokens.json"))
	if err != nil {
		log.Fatalf("No se pudo cargar el repositorio de tokens: %v", err)
	}
	transaccionRepo, err := repositories.NuevoTransaccionRepository(filepath.Join(cfg.DirectorioDatos, "transacciones.json"))
	if err != nil {
		log.Fatalf("No se pudo cargar el repositorio de transacciones: %v", err)
	}

	// Capa de fachada
	suscripcionService := services.NuevoSuscripcionService(suscripcionRepo)
	pagoService := services.NuevoPagoService(tokenRepo, transaccionRepo, suscripcionRepo)

	// Capa de controladores
	router := gin.Default()
	router.Use(middleware.CORS())
	api := router.Group("/api")
	controllers.NuevoSuscripcionController(suscripcionService).Registrar(api)
	controllers.NuevoPagoController(pagoService).Registrar(api)
	controllers.NuevoEstadoController(cfg).Registrar(api)

	log.Println("==============================================")
	log.Println(" Pasarela de pagos de suscripciones Spotify")
	log.Printf(" Puerto:            %s", cfg.Puerto)
	log.Printf(" DIRECTORIO_DATOS:  %s", cfg.DirectorioDatos)
	log.Println("==============================================")

	if err := router.Run(":" + cfg.Puerto); err != nil {
		log.Fatal(err)
	}
}
