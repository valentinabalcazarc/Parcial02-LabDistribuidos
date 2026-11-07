/**
 * @file main.go
 * @brief Punto de entrada principal para el Servidor de Estadísticas.
 *
 * Inicia la conexión con RabbitMQ para procesar eventos asíncronos de reproducción
 * y levanta un servidor HTTP Gin en el puerto :8082 para permitir al cliente consultar
 * los audios que ha reproducido y sus conteos (Requerimiento d).
 */
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	capaaccesoadatos "servidorEstadisticas/capaAccesoADatos"
	capacontrollers "servidorEstadisticas/capaControllers"
	capaservices "servidorEstadisticas/capaServices"
	listener "servidorEstadisticas/componenteListener"
	consumer "servidorEstadisticas/componenteRecibirMensajes"

	"github.com/gin-gonic/gin"

	"servidorEstadisticas/discovery"
	"servidorEstadisticas/middleware"
)

const nombreServicio = "ms-servidorDeEstadisticas"

func main() {
	puerto, _ := strconv.Atoi(env("PORT", "8082"))
	host := env("SERVICE_HOST", "127.0.0.1")
	consulAddr := env("CONSUL_ADDR", "127.0.0.1:8500")
	secreto := env("GATEWAY_SECRET", "secreto-gateway-2026")

	fmt.Println("==================================================")
	fmt.Println("🚀 Iniciando Servidor de Estadísticas...")
	fmt.Println("==================================================")

	// 1. Inicialización de Capas de Datos y Servicio
	repoEstadisticas := capaaccesoadatos.GetRepositorioEstadisticas()
	estadisticasService := capaservices.NewEstadisticasService(repoEstadisticas)
	estadisticasController := capacontrollers.NewEstadisticasController(estadisticasService)

	// 2. Conexión y Consumo de la Cola RabbitMQ (cola_estadisticas)
	rabbitListener, err := listener.NewRabbitListener()
	if err != nil {
		log.Printf("⚠️  [ADVERTENCIA] No se pudo conectar a RabbitMQ: %v\n", err)
		log.Println("   (El servidor continuará para atender consultas HTTP en el puerto :8082)")
	} else {
		defer rabbitListener.Cerrar()
		messageConsumer := consumer.NewMessageConsumer(rabbitListener, estadisticasService)
		err = messageConsumer.IniciarConsumo()
		if err != nil {
			log.Printf("⚠️  [ADVERTENCIA] Error iniciando consumo de mensajes: %v\n", err)
		} else {
			log.Println("✅ Consumidor de RabbitMQ activo y escuchando en 'cola_estadisticas'")
		}
	}

	// 3. Servidor HTTP Gin para Consultas del Cliente
	router := gin.Default()
	router.Use(middleware.FiltroGateway(secreto))

	// Middleware CORS para permitir consultas desde clientes web
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	estadisticasController.RegistrarRutas(router)

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
