package co.edu.unicauca.servidorreacciones.capaConfiguracionWebSocket;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;
import org.springframework.messaging.simp.config.MessageBrokerRegistry;
import org.springframework.web.socket.config.annotation.EnableWebSocketMessageBroker;
import org.springframework.web.socket.config.annotation.StompEndpointRegistry;
import org.springframework.web.socket.config.annotation.WebSocketMessageBrokerConfigurer;

@Configuration
@EnableWebSocketMessageBroker
public class WebSocketConfig implements WebSocketMessageBrokerConfigurer {

  @Value("${reacciones.ws.endpoint}")
  private String endpointWebSocket;

  @Override
  public void configureMessageBroker(MessageBrokerRegistry config) {
    // Los clientes se suscriben a /brokerDeReacciones/audio/{idAudio} para recibir reacciones
    config.enableSimpleBroker("/brokerDeReacciones");
    // Los clientes envían a /app/enviarReaccion
    config.setApplicationDestinationPrefixes("/app");
  }

  @Override
  public void registerStompEndpoints(StompEndpointRegistry registry) {
    registry.addEndpoint(endpointWebSocket)
            .setAllowedOriginPatterns("*")
            .withSockJS();
  }
}
