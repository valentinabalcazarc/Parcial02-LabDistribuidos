package co.edu.unicauca.servidorreacciones.capaListeners;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.event.EventListener;
import org.springframework.messaging.simp.stomp.StompHeaderAccessor;
import org.springframework.stereotype.Component;
import org.springframework.web.socket.messaging.SessionDisconnectEvent;
import org.springframework.web.socket.messaging.SessionSubscribeEvent;
import org.springframework.web.socket.messaging.SessionUnsubscribeEvent;

import co.edu.unicauca.servidorreacciones.capaAccesoADatos.RepositorioReproducciones;
import co.edu.unicauca.servidorreacciones.capaFachadaServices.ServicioReacciones;

/**
 * Sigue qué audio escucha cada sesión (adaptado de PresenciaUsuarioListener de servidorchat):
 * suscripción al topic de un audio, cancelación de la suscripción y desconexión.
 */
@Component
public class SesionWebSocketListener {

    private static final Logger log = LoggerFactory.getLogger(SesionWebSocketListener.class);

    private final RepositorioReproducciones repositorioReproducciones;

    public SesionWebSocketListener(RepositorioReproducciones repositorioReproducciones) {
        this.repositorioReproducciones = repositorioReproducciones;
    }

    // El cliente "empieza a reproducir" un audio cuando se suscribe a su topic
    @EventListener
    public void handleSessionSubscribe(SessionSubscribeEvent event) {
        StompHeaderAccessor accessor = StompHeaderAccessor.wrap(event.getMessage());
        String destino = accessor.getDestination();
        String sessionId = accessor.getSessionId();
        String subscriptionId = accessor.getSubscriptionId();
        if (destino == null || sessionId == null || subscriptionId == null
                || !destino.startsWith(ServicioReacciones.PREFIJO_TOPIC_AUDIO)) {
            return;
        }
        String idAudio = destino.substring(ServicioReacciones.PREFIJO_TOPIC_AUDIO.length());
        repositorioReproducciones.registrar(sessionId, subscriptionId, idAudio);
        log.info("La sesión {} reproduce el audio '{}'", sessionId, idAudio);
    }

    // ... y deja de reproducirlo cuando cancela la suscripción
    @EventListener
    public void handleSessionUnsubscribe(SessionUnsubscribeEvent event) {
        StompHeaderAccessor accessor = StompHeaderAccessor.wrap(event.getMessage());
        String sessionId = accessor.getSessionId();
        String subscriptionId = accessor.getSubscriptionId();
        if (sessionId == null || subscriptionId == null) {
            return;
        }
        repositorioReproducciones.quitarSuscripcion(sessionId, subscriptionId);
    }

    @EventListener
    public void handleSessionDisconnect(SessionDisconnectEvent event) {
        StompHeaderAccessor accessor = StompHeaderAccessor.wrap(event.getMessage());
        String sessionId = accessor.getSessionId();
        if (sessionId == null) {
            return;
        }
        repositorioReproducciones.quitarSesion(sessionId);
        log.info("La sesión {} se ha desconectado", sessionId);
    }
}
