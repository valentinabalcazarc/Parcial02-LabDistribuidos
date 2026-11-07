/**
 * @file message_consumer.go
 * @brief Consumidor de mensajes para procesar notificaciones desde RabbitMQ y persistir estadísticas.
 */
package componenterecibirmensajes

import (
	"encoding/json"
	"fmt"
	"log"

	dtos "servidorEstadisticas/DTOs"
	capaservices "servidorEstadisticas/capaServices"
	listener "servidorEstadisticas/componenteListener"
)

/**
 * @struct MessageConsumer
 * @brief Gestiona la suscripción a la cola y delega el procesamiento al servicio.
 */
type MessageConsumer struct {
	rabbitListener      *listener.RabbitListener
	estadisticasService *capaservices.EstadisticasService
}

/**
 * @brief Crea una nueva instancia de MessageConsumer.
 *
 * @param rabbitListener Listener configurado con la conexión a RabbitMQ.
 * @param service Servicio para persistir y procesar las estadísticas.
 */
func NewMessageConsumer(rabbitListener *listener.RabbitListener, service *capaservices.EstadisticasService) *MessageConsumer {
	return &MessageConsumer{
		rabbitListener:      rabbitListener,
		estadisticasService: service,
	}
}

/**
 * @brief Inicia el consumo continuo de mensajes desde la cola de RabbitMQ.
 */
func (mc *MessageConsumer) IniciarConsumo() error {
	ch := mc.rabbitListener.GetChannel()
	queueName := mc.rabbitListener.GetQueueName()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("Error al registrar el consumidor: %v", err)
	}

	log.Println(" [*] Esperando mensajes de reproducción en 'cola_estadisticas'...")

	go func() {
		for d := range msgs {
			var notificacion dtos.NotificacionReproduccion
			err := json.Unmarshal(d.Body, &notificacion)
			if err != nil {
				log.Printf("Error decodificando mensaje JSON: %v", err)
				continue
			}

			// 1. Guardar en el servicio/repositorio de estadísticas
			mc.estadisticasService.ProcesarNotificacion(notificacion)

			// 2. Imprimir por consola con formato descriptivo
			fmt.Println("\n==================================================")
			fmt.Println("   📊 [SERVIDOR DE ESTADÍSTICAS - EVENTO PROCESADO]")
			if notificacion.Cliente != "" {
				fmt.Printf("   👤 Cliente:    %s\n", notificacion.Cliente)
			}
			fmt.Printf("   🎵 Audio:      %s\n", notificacion.Titulo)
			if notificacion.TipoAudio != "" {
				fmt.Printf("   🏷️  Tipo:       %s\n", notificacion.TipoAudio)
			}
			fmt.Printf("   🕒 Fecha/Hora: %s\n", notificacion.FechaHora)
			fmt.Printf("   💬 Mensaje:    %s\n", notificacion.Mensaje)
			fmt.Println("==================================================")
		}
	}()

	return nil
}
