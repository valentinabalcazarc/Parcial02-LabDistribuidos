package co.edu.unicauca.servidorreacciones.capaControladores;

import org.springframework.messaging.handler.annotation.MessageMapping;
import org.springframework.messaging.simp.SimpMessageHeaderAccessor;
import org.springframework.stereotype.Controller;

import co.edu.unicauca.servidorreacciones.capaFachadaServices.ServicioReacciones;
import co.edu.unicauca.servidorreacciones.capaFachadaServices.DTO.ReaccionDTO;

@Controller
public class ReaccionesController {

    private final ServicioReacciones servicioReacciones;

    public ReaccionesController(ServicioReacciones servicioReacciones) {
        this.servicioReacciones = servicioReacciones;
    }

    // El cliente envía a /app/enviarReaccion; el servicio la reenvía a los que reproducen el mismo audio
    @MessageMapping("/enviarReaccion")
    public void enviarReaccion(ReaccionDTO reaccion, SimpMessageHeaderAccessor headerAccessor) {
        servicioReacciones.enviarReaccion(headerAccessor.getSessionId(), reaccion);
    }

}
