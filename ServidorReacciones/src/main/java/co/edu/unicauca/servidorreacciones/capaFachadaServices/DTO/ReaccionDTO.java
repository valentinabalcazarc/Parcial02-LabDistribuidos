package co.edu.unicauca.servidorreacciones.capaFachadaServices.DTO;

import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * Reacción asociada a un audio. El servidor la reenvía tal cual a los que escuchan ese audio.
 */
@Data
@NoArgsConstructor
public class ReaccionDTO {
    private String idAudio;
    private String tipoReaccion;
}
