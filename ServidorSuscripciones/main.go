package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"pasarela-pagos/config"
	"pasarela-pagos/controllers"
	"pasarela-pagos/data/repositories"
	"pasarela-pagos/discovery"
	"pasarela-pagos/middleware"
	"pasarela-pagos/services"
)

func main() {
	cfg := config.Cargar()
	secreto := env("GATEWAY_SECRET", "secreto-gateway-2026")

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
	router.Use(middleware.FiltroGateway(secreto))
	api := router.Group("/api")
	controllers.NuevoSuscripcionController(suscripcionService).Registrar(api)
	controllers.NuevoPagoController(pagoService).Registrar(api)
	controllers.NuevoEstadoController(cfg).Registrar(api)

	log.Println("==============================================")
	log.Println(" Pasarela de pagos de suscripciones Spotify")
	log.Printf(" Puerto:            %s", cfg.Puerto)
	log.Printf(" DIRECTORIO_DATOS:  %s", cfg.DirectorioDatos)
	log.Printf(" Consul:            %s", cfg.ConsulAddr)
	log.Printf(" Nombre servicio:   %s", cfg.NombreServicio)
	log.Println("==============================================")

	puerto, _ := strconv.Atoi(cfg.Puerto)
	srv := &http.Server{Addr: fmt.Sprintf(":%d", puerto), Handler: router}
	go func() {
		log.Printf("%s escuchando en :%d", cfg.NombreServicio, puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error en el servidor: %v", err)
		}
	}()

	// Registro en Consul para que el API Gateway descubra este servicio.
	registro, err := discovery.Registrar(cfg.ConsulAddr, cfg.NombreServicio, cfg.Host, puerto)
	if err != nil {
		log.Fatalf("no fue posible registrarse en Consul: %v", err)
	}

	// Apagado limpio: desregistrar de Consul y cerrar el servidor HTTP.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	registro.Desregistrar()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func env(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}
