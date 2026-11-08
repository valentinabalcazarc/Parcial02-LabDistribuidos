package co.edu.unicauca.servidorreacciones.capaFachadaServices;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.stereotype.Service;

import co.edu.unicauca.servidorreacciones.capaAccesoADatos.RepositorioReproducciones;
import co.edu.unicauca.servidorreacciones.capaFachadaServices.DTO.ReaccionDTO;

/**
 * Lógica de las reacciones: valida que el emisor esté escuchando el audio
 * y reenvía la reacción a todos los que escuchan ese mismo audio.
 */
@Service
public class ServicioReacciones {

    private static final Logger log = LoggerFactory.getLogger(ServicioReacciones.class);

    /** Topic de cada audio: /brokerDeReacciones/audio/{idAudio} */
    public static final String PREFIJO_TOPIC_AUDIO = "/brokerDeReacciones/audio/";

    private static final int LONGITUD_MAXIMA_TIPO = 50;

    private final SimpMessagingTemplate simpMessagingTemplate;
    private final RepositorioReproducciones repositorioReproducciones;

    public ServicioReacciones(SimpMessagingTemplate simpMessagingTemplate,
                              RepositorioReproducciones repositorioReproducciones) {
        this.simpMessagingTemplate = simpMessagingTemplate;
        this.repositorioReproducciones = repositorioReproducciones;
    }

    /**
     * @param sessionId sesión WebSocket que envía la reacción
     * @param reaccion  reacción recibida (idAudio y tipoReaccion)
     * @return true si la reacción fue reenviada a los oyentes del audio
     */
    public boolean enviarReaccion(String sessionId, ReaccionDTO reaccion) {
        if (!esValida(reaccion)) {
            log.warn("Reacción descartada (datos incompletos o inválidos) de la sesión {}", sessionId);
            return false;
        }

        String idAudio = reaccion.getIdAudio();
        if (!repositorioReproducciones.estaReproduciendo(sessionId, idAudio)) {
            log.warn("Reacción descartada: la sesión {} no está reproduciendo '{}'", sessionId, idAudio);
            return false;
        }

        // Callback: cada cliente suscrito a este topic recibe la reacción por su conexión ya abierta
        simpMessagingTemplate.convertAndSend(PREFIJO_TOPIC_AUDIO + idAudio, reaccion);
        return true;
    }

    private boolean esValida(ReaccionDTO reaccion) {
        return reaccion != null
                && reaccion.getIdAudio() != null && !reaccion.getIdAudio().isBlank()
                && reaccion.getTipoReaccion() != null && !reaccion.getTipoReaccion().isBlank()
                && reaccion.getTipoReaccion().length() <= LONGITUD_MAXIMA_TIPO;
    }
}
