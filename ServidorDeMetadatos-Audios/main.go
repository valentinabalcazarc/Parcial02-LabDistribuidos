package main

import (
	repository "servidorMetadatos/capaAccesoADatos/repository"
	controller "servidorMetadatos/capaControllers"
	service "servidorMetadatos/capaServices/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// TIPOS DE AUDIO
	tipoAudioRepository := repository.NewTipoAudioRepository()
	tipoAudioService := service.NewTipoAudioService(tipoAudioRepository)
	tipoAudioController := controller.NewTipoAudioController(tipoAudioService)

	// AUDIOLIBRO
	audiolibroRepository := repository.NewMetadataAudiolibroRepository()
	audiolibroService := service.NewMetadataAudiolibroService(audiolibroRepository)
	audiolibroController := controller.NewMetadataAudiolibroController(audiolibroService)

	// MUSICA
	musicaRepository := repository.NewMetadataMusicaRepository()
	musicaService := service.NewMetadataMusicaService(musicaRepository)
	musicaController := controller.NewMetadataMusicaController(musicaService)

	// PODCAST
	podcastRepository := repository.NewMetadataPodcastRepository()
	podcastService := service.NewMetadataPodcastService(podcastRepository)
	podcastController := controller.NewMetadataPodcastController(podcastService)

	// RUIDO BLANCO
	ruidoBlancoRepository := repository.NewMetadataRuidoBlancoRepository()
	ruidoBlancoService := service.NewMetadataRuidoBlancoService(ruidoBlancoRepository)
	ruidoBlancoController := controller.NewMetadataRuidoBlancoController(ruidoBlancoService)

	// ROUTER GIN
	router := gin.Default()

	// Ruta para obtener todos los tipos de audio
	router.GET("/tipos-audio", tipoAudioController.ListarTipos)

	// Rutas de consulta por ID (en singular)
	router.GET("/audiolibro/:id", audiolibroController.ConsultarAudiolibro)
	router.GET("/musica/:id", musicaController.ConsultarMusica)
	router.GET("/podcast/:id", podcastController.ConsultarPodcast)
	router.GET("/ruidoblanco/:id", ruidoBlancoController.ConsultarRuidoBlanco)

	// Rutas de listado completo (en plural)
	router.GET("/audiolibros", audiolibroController.ListarAudiolibros)
	router.GET("/musicas", musicaController.ListarMusica)
	router.GET("/podcasts", podcastController.ListarPodcasts)
	router.GET("/ruidosblancos", ruidoBlancoController.ListarRuidoBlanco)

	// Iniciar servidor en el puerto 8081
	router.Run(":8081")
}
