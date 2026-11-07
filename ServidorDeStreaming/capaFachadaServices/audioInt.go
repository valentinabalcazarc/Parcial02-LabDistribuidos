package capafachadaservices

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	capaAccesoDatos "servidor.local/grpc-servidor/capaAccesoADatos"
	componenteconexioncola "servidor.local/grpc-servidor/componenteConexionCola"
)

type StreamingService struct {
	repo         *capaAccesoDatos.RepositorioStreaming
	conexionCola *componenteconexioncola.RabbitPublisher
}

/**
 * @brief Crea una nueva instancia de StreamingService.
 * @details Inicializa el repositorio y la conexión a la cola de RabbitMQ.
 * @return *StreamingService Puntero a la nueva instancia del servicio.
 */
func NewStreamingService() *StreamingService {
	fmt.Println("Inicializando service de streaming")
	repo := capaAccesoDatos.GetRepositorioStreaming()
	conexionCola, err := componenteconexioncola.NewRabbitPublisher()

	if err != nil {
		fmt.Println("Error al conectar con RabbitMQ: ", err)
		conexionCola = nil
	}

	return &StreamingService{
		repo:         repo,
		conexionCola: conexionCola,
	}
}

/**
 * @brief Actúa como fachada para obtener el archivo de la canción.
 * @details Usa la capa de acceso a datos para abrir el archivo.
 * @param titulo Nombre o título del archivo de audio a abrir.
 * @return *os.File Puntero al archivo abierto.
 * @return error Error si el archivo no se puede abrir.
 */
func abrirArchivo(titulo string) (*os.File, error) {
	log.Printf("GetAudioFile llamado con titulo=%s", titulo)
	return capaAccesoDatos.GetRepositorioStreaming().AbrirArchivo(titulo)
}

/**
 * @brief Lee el archivo de audio y lo transmite en bloques de 64 KB, notificando a RabbitMQ.
 *
 * @param tituloCancion Nombre del archivo a transmitir.
 * @param cliente Identificador del cliente que solicita la reproducción (para estadísticas).
 * @param funcionParaEnviarFragmento Callback de streaming gRPC.
 */
func StreamAudioFile(tituloCancion string, cliente string, funcionParaEnviarFragmento func([]byte) error) error {
	log.Printf("Canción solicitada: %s por cliente: %s", tituloCancion, cliente)
	file, err := abrirArchivo(tituloCancion)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer file.Close()

	// 1. Notificar inicio de reproducción a RabbitMQ para estadísticas
	pub, err := componenteconexioncola.NewRabbitPublisher()
	if err == nil {
		defer pub.Cerrar()
		msg := componenteconexioncola.NotificacionReproduccion{
			Cliente:   cliente,
			Titulo:    tituloCancion,
			FechaHora: time.Now().Format("2006-01-02 15:04:05"),
			Mensaje:   fmt.Sprintf("Reproduciendo audio %s por streaming", tituloCancion),
		}
		if errPub := pub.PublicarNotificacion(msg); errPub != nil {
			log.Printf("Advertencia al publicar en RabbitMQ: %v", errPub)
		}
	} else {
		log.Printf("Advertencia: No se pudo conectar a RabbitMQ para enviar notificación: %v", err)
	}

	// 2. Transmisión del archivo en bloques de 64 KB
	buffer := make([]byte, 64*1024)
	fragmento := 0

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			log.Println("Canción enviada completamente desde la fachada.")
			break
		}
		if err != nil {
			return fmt.Errorf("error leyendo el archivo: %w", err)
		}

		fragmento++
		log.Printf("Fragmento #%d leído (%d bytes) y enviando", fragmento, n)
		err = funcionParaEnviarFragmento(buffer[:n])
		if err != nil {
			return fmt.Errorf("error enviando fragmento #%d: %w", fragmento, err)
		}
	}

	return nil
}
