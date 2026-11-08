package co.edu.unicauca.servidorreacciones.capaAccesoADatos;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import org.springframework.stereotype.Repository;

/**
 * Qué audio está reproduciendo cada sesión.
 * Un cliente "reproduce" un audio mientras mantenga una suscripción a su topic,
 * por eso se guarda sessionId -> (subscriptionId -> idAudio).
 */
@Repository
public class RepositorioReproducciones {

    private final ConcurrentHashMap<String, Map<String, String>> audiosPorSesion = new ConcurrentHashMap<>();

    public void registrar(String sessionId, String subscriptionId, String idAudio) {
        audiosPorSesion
                .computeIfAbsent(sessionId, k -> new ConcurrentHashMap<>())
                .put(subscriptionId, idAudio);
    }

    public void quitarSuscripcion(String sessionId, String subscriptionId) {
        Map<String, String> suscripciones = audiosPorSesion.get(sessionId);
        if (suscripciones != null) {
            suscripciones.remove(subscriptionId);
        }
    }

    public void quitarSesion(String sessionId) {
        audiosPorSesion.remove(sessionId);
    }

    public boolean estaReproduciendo(String sessionId, String idAudio) {
        Map<String, String> suscripciones = audiosPorSesion.get(sessionId);
        return suscripciones != null && suscripciones.containsValue(idAudio);
    }
}
