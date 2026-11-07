package main

import (
	repository "servidorMetadatos/capaAccesoADatos/repository"
	controller "servidorMetadatos/capaControllers"
	service "servidorMetadatos/capaServices/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// ==========================================
	// 1. INICIALIZACIÓN DE SERVICIO DE AUDIOS FÍSICOS
	// ==========================================
	// Este servicio centraliza el guardado de archivos binarios (.mp3)
	// para que ningún controlador repita lógica de disco.
	audioService := service.NewAudioService()
	audioController := controller.NewAudioController(audioService)

	// ==========================================
	// 2. TIPOS DE AUDIO
	// ==========================================
	tipoAudioRepository := repository.NewTipoAudioRepository()
	tipoAudioService := service.NewTipoAudioService(tipoAudioRepository)
	tipoAudioController := controller.NewTipoAudioController(tipoAudioService)

	// ==========================================
	// 3. AUDIOLIBRO
	// ==========================================
	audiolibroRepository := repository.NewMetadataAudiolibroRepository()
	audiolibroService := service.NewMetadataAudiolibroService(audiolibroRepository, audioService)
	audiolibroController := controller.NewMetadataAudiolibroController(audiolibroService)

	// ==========================================
	// 4. MÚSICA
	// ==========================================
	musicaRepository := repository.NewMetadataMusicaRepository()
	musicaService := service.NewMetadataMusicaService(musicaRepository, audioService)
	musicaController := controller.NewMetadataMusicaController(musicaService)

	// ==========================================
	// 5. PODCAST
	// ==========================================
	podcastRepository := repository.NewMetadataPodcastRepository()
	podcastService := service.NewMetadataPodcastService(podcastRepository, audioService)
	podcastController := controller.NewMetadataPodcastController(podcastService)

	// ==========================================
	// 6. RUIDO BLANCO
	// ==========================================
	ruidoBlancoRepository := repository.NewMetadataRuidoBlancoRepository()
	ruidoBlancoService := service.NewMetadataRuidoBlancoService(ruidoBlancoRepository, audioService)
	ruidoBlancoController := controller.NewMetadataRuidoBlancoController(ruidoBlancoService)

	// ==========================================
	// 7. ROUTER GIN Y MIDDLEWARES
	// ==========================================
	router := gin.Default()

	// Middleware CORS para habilitar peticiones desde el cliente web
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// --- Rutas generales ---
	router.GET("/tipos-audio", tipoAudioController.ListarTipos)

	// --- Rutas de consulta por ID (GET en singular) ---
	router.GET("/audiolibro/:id", audiolibroController.ConsultarAudiolibro)
	router.GET("/musica/:id", musicaController.ConsultarMusica)
	router.GET("/podcast/:id", podcastController.ConsultarPodcast)
	router.GET("/ruidoblanco/:id", ruidoBlancoController.ConsultarRuidoBlanco)

	// --- Rutas de listado completo (GET en plural) ---
	router.GET("/audiolibros", audiolibroController.ListarAudiolibros)
	router.GET("/musicas", musicaController.ListarMusica)
	router.GET("/podcasts", podcastController.ListarPodcasts)
	router.GET("/ruidosblancos", ruidoBlancoController.ListarRuidoBlanco)

	// --- Rutas de registro unificado (POST: Metadatos + Archivo físico) ---
	// La llamada a metadatos coordina internamente el guardado del audio
	router.POST("/musica", musicaController.RegistrarMusica)
	router.POST("/audiolibro", audiolibroController.RegistrarAudiolibro)
	router.POST("/podcast", podcastController.RegistrarPodcast)
	router.POST("/ruidoblanco", ruidoBlancoController.RegistrarRuidoBlanco)

	// --- Ruta independiente para subir solo archivos de audio físicos ---
	router.POST("/audio", audioController.GuardarAudio)

	// Iniciar servidor en el puerto 8081
	router.Run(":8081")
}
