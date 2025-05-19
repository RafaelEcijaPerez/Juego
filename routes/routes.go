package routes

import (
	"github.com/gin-gonic/gin"
	"juego/controllers"
	"juego/db"
	"juego/handlers"
	"juego/services"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Crear una instancia del JugadorService
	// Asegúrate de que db.DB sea del tipo *gorm.io/gorm.DB
	jugadorService := services.NewJugadorService(db.DB)

	// Nuevo controlador QR
	qrController := controllers.NewQRController()

	// Rutas QR
	r.POST("/generate-qr", qrController.GenerateQR)
	r.POST("/check-qr-status", qrController.CheckQRStatus)
	r.POST("/login/qr", qrController.LoginWithQR)
	// Crear una instancia del JugadorController
	jugadorController := controllers.NewJugadorController(jugadorService)

	// Ruta para el registro
	r.POST("/register", jugadorController.Register)

	// Ruta para iniciar sesión
	r.POST("/login", jugadorController.Login)

	// Rutas para el juego Cuatro en Raya
	r.POST("/crear-cuatro-en-raya", handlers.CrearJuego)
	r.GET("/obtener-cuatro-en-raya/:id", handlers.ObtenerJuego)
	r.POST("/movimiento-cuatro-en-raya/:id", handlers.HacerMovimiento)
	r.POST("/terminar-cuatro-en-raya/:id", handlers.TerminarJuego)

	// Rutas para el juego conecta Cuatro
	r.POST("/crear-conecta-cuatro", handlers.CrearJuegoConecta)
	r.GET("/obtener-conecta-cuatro/:id", handlers.ObtenerJuegoConecta)
	r.POST("/movimiento-conecta-cuatro/:id", handlers.HacerMovimientoConecta)
	r.POST("/terminar-conecta-cuatro/:id", handlers.TerminarJuegoConecta)

	// Rutas para el juego Desde el borde
	r.POST("/crear-desde-borde", handlers.CrearJuegoDesdeBorde)
	r.GET("/obtener-desde-borde/:id", handlers.ObtenerJuegoDesdeBorde)
	r.POST("/movimiento-desde-borde/:id", handlers.HacerMovimientoDesdeBorde)
	r.POST("/terminar-desde-borde/:id", handlers.TerminarJuegoDesdeBorde)

	// Rutas para el juego Pasa Bolas
	r.POST("/crear-juego-pasa-bolas", handlers.CrearJuegoPasaBolas)
	r.GET("/obtener-juego-pasa-bolas/:id", handlers.ObtenerJuegoPasaBolas)
	r.POST("/lanza-bola-pasa-bolas/:id", handlers.LanzarBola)
	r.POST("/terminar-juego-pasa-bolas/:id", handlers.TerminarJuegoPasaBolas)
	r.POST("/reiniciar-juego-pasa-bolas/:id", handlers.ReiniciarBolasPasaBolas)

	return r
}
