/**
 * @file rabbitmq_publisher.go
 * @brief Componente para la conexión y publicación de mensajes en RabbitMQ.
 */
package componenteconexioncola

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

/**
 * @struct RabbitPublisher
 * @brief Estructura que maneja la conexión, el canal y la cola de RabbitMQ.
 */
type RabbitPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

/**
 * @struct NotificacionReproduccion
 * @brief Estructura que representa el mensaje de notificación de reproducción a enviar.
 */
type NotificacionReproduccion struct {
	Cliente   string `json:"cliente,omitempty"`    /**< @brief Identificador del cliente que reproduce. */
	Titulo    string `json:"titulo"`               /**< @brief Título del audio reproducido. */
	TipoAudio string `json:"tipo_audio,omitempty"` /**< @brief Categoría de audio opcional. */
	FechaHora string `json:"fecha_hora,omitempty"` /**< @brief Fecha y hora del evento. */
	Mensaje   string `json:"mensaje"`              /**< @brief Mensaje descriptivo. */
}

/**
 * @brief Crea una nueva conexión y publicador de RabbitMQ.
 * @details Configura la conexión, abre un canal y declara la cola "cola_estadisticas".
 * @return *RabbitPublisher Puntero al publicador creado.
 * @return error Error en caso de fallo en la conexión o configuración.
 */
func NewRabbitPublisher() (*RabbitPublisher, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://admin:1234@172.20.10.2:5672/"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("Error conectando a rabbitMQ en %s: %v", url, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Error abriendo el canal: %v", err)
	}

	q, err := ch.QueueDeclare(
		"cola_estadisticas", // nombre de la cola
		true,                // durable
		false,               // autodelete
		false,               // exclusive
		false,               // no-wait
		nil,                 // args
	)

	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("Error declarando la cola: %v", err)
	}

	return &RabbitPublisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

/**
 * @brief Publica una notificación de reproducción en la cola de RabbitMQ.
 * @param msg Mensaje de notificación a publicar.
 * @return error Error en caso de que falle la conversión a JSON o la publicación.
 */
func (p *RabbitPublisher) PublicarNotificacion(msg NotificacionReproduccion) error {
	body, err := json.Marshal(msg)

	if err != nil {
		return fmt.Errorf("Error convirtiendo a JSON: %v", err)
	}

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("Error publicando el mensaje: %v", err)
	}

	fmt.Println("[RabbitMQ] Notificación enviada a la cola:", string(body))
	return nil
}

/**
 * @brief Cierra el canal y la conexión con RabbitMQ.
 */
func (p *RabbitPublisher) Cerrar() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
