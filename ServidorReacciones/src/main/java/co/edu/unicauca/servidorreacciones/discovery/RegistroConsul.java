package co.edu.unicauca.servidorreacciones.discovery;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;

import com.fasterxml.jackson.databind.ObjectMapper;

import jakarta.annotation.PreDestroy;

/**
 * Registro del microservicio en Consul con un health check TTL que el propio servicio renueva
 * (mismo esquema que discovery/consul.go de streaming y metadatos: TTL 15s, latido cada 5s).
 * Si Consul no está disponible, el servicio sigue funcionando y reintenta el registro en cada latido.
 */
@Component
public class RegistroConsul {

    private static final Logger log = LoggerFactory.getLogger(RegistroConsul.class);

    private final HttpClient http = HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(3)).build();
    private final ObjectMapper mapper = new ObjectMapper();
    private final ScheduledExecutorService latido = Executors.newSingleThreadScheduledExecutor(tarea -> {
        Thread hilo = new Thread(tarea, "consul-latido");
        hilo.setDaemon(true);
        return hilo;
    });

    private final String consulAddr;
    private final String nombre;
    private final String host;
    private final int puerto;
    private final String id;
    private final String checkId;
    private volatile boolean registrado = false;

    public RegistroConsul(@Value("${consul.addr}") String consulAddr,
                          @Value("${servicio.nombre}") String nombre,
                          @Value("${servicio.host}") String host,
                          @Value("${server.port}") int puerto) {
        this.consulAddr = consulAddr;
        this.nombre = nombre;
        this.host = host;
        this.puerto = puerto;
        this.id = nombre + "-" + host + "-" + puerto;
        this.checkId = "check-" + id;
    }

    @EventListener(ApplicationReadyEvent.class)
    public void iniciar() {
        latido.scheduleAtFixedRate(this::latir, 0, 5, TimeUnit.SECONDS);
    }

    private void latir() {
        try {
            if (!registrado) {
                registrar();
            }
            enviarPut("/v1/agent/check/update/" + checkId, "{\"Status\":\"passing\",\"Output\":\"ok\"}");
        } catch (Exception e) {
            if (e instanceof InterruptedException) {
                Thread.currentThread().interrupt();
            }
            registrado = false;
            String motivo = e.getMessage() != null ? e.getMessage() : e.getClass().getSimpleName();
            log.warn("Error comunicándose con Consul ({}): {}", consulAddr, motivo);
        }
    }

    private void registrar() throws Exception {
        Map<String, Object> check = new LinkedHashMap<>();
        check.put("CheckID", checkId);
        check.put("TTL", "15s");
        check.put("DeregisterCriticalServiceAfter", "1m");

        Map<String, Object> servicio = new LinkedHashMap<>();
        servicio.put("ID", id);
        servicio.put("Name", nombre);
        servicio.put("Address", host);
        servicio.put("Port", puerto);
        servicio.put("Tags", List.of("websocket", "stomp"));
        servicio.put("Check", check);

        enviarPut("/v1/agent/service/register", mapper.writeValueAsString(servicio));
        registrado = true;
        log.info("Servicio '{}' registrado en Consul ({}) como {}", nombre, consulAddr, id);
    }

    // Al apagarse, el servicio se elimina del directorio de descubrimiento
    @PreDestroy
    public void desregistrar() {
        latido.shutdownNow();
        if (!registrado) {
            return;
        }
        try {
            enviarPut("/v1/agent/service/deregister/" + id, null);
            log.info("Servicio {} desregistrado de Consul", id);
        } catch (Exception e) {
            log.warn("Error desregistrando de Consul: {}", e.getMessage());
        }
    }

    private void enviarPut(String ruta, String cuerpo) throws IOException, InterruptedException {
        HttpRequest.BodyPublisher publicador = cuerpo == null
                ? HttpRequest.BodyPublishers.noBody()
                : HttpRequest.BodyPublishers.ofString(cuerpo);
        HttpRequest peticion = HttpRequest.newBuilder(URI.create("http://" + consulAddr + ruta))
                .timeout(Duration.ofSeconds(3))
                .header("Content-Type", "application/json")
                .PUT(publicador)
                .build();
        HttpResponse<String> respuesta = http.send(peticion, HttpResponse.BodyHandlers.ofString());
        if (respuesta.statusCode() / 100 != 2) {
            throw new IOException("Consul respondió " + respuesta.statusCode() + ": " + respuesta.body());
        }
    }
}
