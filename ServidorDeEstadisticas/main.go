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
	"fmt"
	"log"

	capaaccesoadatos "servidorEstadisticas/capaAccesoADatos"
	capacontrollers "servidorEstadisticas/capaControllers"
	capaservices "servidorEstadisticas/capaServices"
	listener "servidorEstadisticas/componenteListener"
	consumer "servidorEstadisticas/componenteRecibirMensajes"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// 3. Servidor HTTP Gin para Consultas del Cliente (Requerimiento d)
	router := gin.Default()

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

	// --- Rutas de Consulta de Estadísticas ---
	// Requerimiento (d): Un cliente puede ver los audios que ha reproducido y cuántas veces
	router.GET("/estadisticas/cliente/:cliente", estadisticasController.ConsultarReproduccionesCliente)

	// Estadísticas globales de popularidad de audios
	router.GET("/estadisticas/audios", estadisticasController.ConsultarRankingGlobal)

	// Historial de eventos recientes
	router.GET("/estadisticas/historial", estadisticasController.ConsultarHistorial)

	fmt.Println("📡 Servidor HTTP de Estadísticas escuchando en el puerto :8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("Error fatal iniciando servidor HTTP en :8082: %v", err)
	}
}
