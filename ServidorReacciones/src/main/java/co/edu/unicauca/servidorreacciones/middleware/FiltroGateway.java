package co.edu.unicauca.servidorreacciones.middleware;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

/**
 * Rechaza toda petición que no traiga el token secreto del API Gateway
 * (equivalente a FiltroGateway de los servidores en Go). Cubre el handshake
 * del WebSocket y las rutas de SockJS.
 */
@Component
@ConditionalOnProperty(name = "gateway.habilitado", havingValue = "true", matchIfMissing = true)
public class FiltroGateway extends OncePerRequestFilter {

    /** Cabecera que solo el API Gateway agrega a las peticiones. */
    public static final String HEADER_GATEWAY = "X-Gateway-Token";

    private final byte[] secreto;

    public FiltroGateway(@Value("${gateway.secreto}") String secreto) {
        this.secreto = secreto.getBytes(StandardCharsets.UTF_8);
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain chain)
            throws ServletException, IOException {
        String token = request.getHeader(HEADER_GATEWAY);
        // isEqual compara en tiempo constante (como subtle.ConstantTimeCompare en Go)
        if (token == null || token.isEmpty()
                || !MessageDigest.isEqual(token.getBytes(StandardCharsets.UTF_8), secreto)) {
            response.setStatus(HttpServletResponse.SC_FORBIDDEN);
            response.setContentType("application/json;charset=UTF-8");
            response.getWriter().write("{\"error\":\"acceso denegado: la petición debe pasar por el API Gateway\"}");
            return;
        }
        chain.doFilter(request, response);
    }
}
