package routes

import (
	"juego/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Rutas para el juego Cuatro en Raya
	r.POST("/crear-juego-cuatro-en-raya", handlers.CrearJuego)
	r.GET("/obtener-juego-cuatro-en-raya/:id", handlers.ObtenerJuego)
	r.POST("/hacer-movimiento-cuatro-en-raya/:id", handlers.HacerMovimiento)
	r.POST("/terminar-juego-cuatro-en-raya/:id", handlers.TerminarJuego)

	// Rutas para el juego Pasa Bolas
	r.POST("/crear-juego-pasa-bolas", handlers.CrearJuegoPasaBolas)
	r.GET("/obtener-juego-pasa-bolas/:id", handlers.ObtenerJuegoPasaBolas)
	r.POST("/eliminar-jugador-pasa-bolas/:id", handlers.EliminarJugadorPasaBolas)
	r.POST("/terminar-juego-pasa-bolas/:id", handlers.TerminarJuegoPasaBolas)

	return r
}