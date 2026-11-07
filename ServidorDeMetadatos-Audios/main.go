package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	repository "servidorMetadatos/capaAccesoADatos/repository"
	controller "servidorMetadatos/capaControllers"
	service "servidorMetadatos/capaServices/services"
	"strconv"
	"syscall"
	"time"

	"servidorMetadatos/discovery"
	"servidorMetadatos/middleware"

	"github.com/gin-gonic/gin"
)

const nombreServicio = "ms-servidorDeMetadatos-Audios"

func main() {
	puerto, _ := strconv.Atoi(env("PORT", "8081"))
	host := env("SERVICE_HOST", "127.0.0.1")
	consulAddr := env("CONSUL_ADDR", "127.0.0.1:8500")
	secreto := env("GATEWAY_SECRET", "secreto-gateway-2026")

	audioService := service.NewAudioService()
	audioController := controller.NewAudioController(audioService)

	tipoAudioRepository := repository.NewTipoAudioRepository()
	tipoAudioService := service.NewTipoAudioService(tipoAudioRepository)
	tipoAudioController := controller.NewTipoAudioController(tipoAudioService)

	audiolibroRepository := repository.NewMetadataAudiolibroRepository()
	audiolibroService := service.NewMetadataAudiolibroService(audiolibroRepository, audioService)
	audiolibroController := controller.NewMetadataAudiolibroController(audiolibroService)

	musicaRepository := repository.NewMetadataMusicaRepository()
	musicaService := service.NewMetadataMusicaService(musicaRepository, audioService)
	musicaController := controller.NewMetadataMusicaController(musicaService)

	podcastRepository := repository.NewMetadataPodcastRepository()
	podcastService := service.NewMetadataPodcastService(podcastRepository, audioService)
	podcastController := controller.NewMetadataPodcastController(podcastService)

	ruidoBlancoRepository := repository.NewMetadataRuidoBlancoRepository()
	ruidoBlancoService := service.NewMetadataRuidoBlancoService(ruidoBlancoRepository, audioService)
	ruidoBlancoController := controller.NewMetadataRuidoBlancoController(ruidoBlancoService)

	router := gin.Default()
	router.Use(middleware.FiltroGateway(secreto))

	// Middleware CORS para habilitar peticiones desde el cliente web
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	audioController.RegistrarRutas(router)
	tipoAudioController.RegistrarRutas(router)
	audiolibroController.RegistrarRutas(router)
	musicaController.RegistrarRutas(router)
	podcastController.RegistrarRutas(router)
	ruidoBlancoController.RegistrarRutas(router)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", puerto), Handler: router}
	go func() {
		log.Printf("%s escuchando en :%d", nombreServicio, puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error en el servidor: %v", err)
		}
	}()

	// Al lanzarse, el microservicio se registra en el directorio de descubrimiento.
	registro, err := discovery.Registrar(consulAddr, nombreServicio, host, puerto)
	if err != nil {
		log.Fatalf("no fue posible registrarse en Consul: %v", err)
	}

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
