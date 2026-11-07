/**
 * @file rabbitmq_listener_config.go
 * @brief Configuración y gestión de la conexión con RabbitMQ.
 */
package componentelistener

import (
	"fmt"

	"github.com/streadway/amqp"
)

/**
 * @struct RabbitListener
 * @brief Estructura que mantiene el estado de la conexión a RabbitMQ.
 */
type RabbitListener struct {
	conn    *amqp.Connection /**< @brief Conexión principal al servidor RabbitMQ. */
	channel *amqp.Channel    /**< @brief Canal de comunicación sobre la conexión. */
	queue   amqp.Queue       /**< @brief Cola declarada para consumir mensajes. */
}

/**
 * @brief Crea una nueva instancia de RabbitListener.
 *
 * Establece la conexión con el servidor RabbitMQ, abre un canal y declara la cola
 * "cola_estadisticas".
 *
 * @return Un puntero a RabbitListener o error en caso de fallo.
 */
func NewRabbitListener() (*RabbitListener, error) {
	// Mismos datos de conexión usados en ServidorDeStreaming
	//conn, err := amqp.Dial("amqp://admin:1234@192.168.80.25:5672/")
	conn, err := amqp.Dial("amqp://admin:1234@172.20.10.2:5672/")

	if err != nil {
		return nil, fmt.Errorf("Error conectando a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Error abriendo el canal: %v", err)
	}

	// Declarar exactamente la misma cola que el publicador
	q, err := ch.QueueDeclare(
		"cola_estadisticas", // nombre de la cola
		true,                // durable
		false,               // autodelete
		false,               // exclusive
		false,               // no-wait
		nil,                 // args
	)

	if err != nil {
		return nil, fmt.Errorf("Error declarando la cola: %v", err)
	}

	fmt.Printf("Conectado exitosamente a RabbitMQ\n")

	return &RabbitListener{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

/**
 * @brief Obtiene el canal de comunicación actual.
 *
 * @return El canal de RabbitMQ asociado al listener.
 */
func (l *RabbitListener) GetChannel() *amqp.Channel {
	return l.channel
}

/**
 * @brief Obtiene el nombre de la cola configurada.
 *
 * @return El nombre de la cola en formato string.
 */
func (l *RabbitListener) GetQueueName() string {
	return l.queue.Name
}

/**
 * @brief Cierra de forma segura el canal y la conexión a RabbitMQ.
 */
func (l *RabbitListener) Cerrar() {
	if l.channel != nil {
		l.channel.Close()
	}
	if l.conn != nil {
		l.conn.Close()
	}
}
