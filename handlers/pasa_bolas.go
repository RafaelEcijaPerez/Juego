package handlers

import (
	"fmt"
	"juego/models"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	juegosActivosPasaBolas = make(map[string]models.PasaBolas) // Juegos activos
	juegosMutexPasaBolas   sync.RWMutex                        // Mutex para acceder a los datos del juego
)

// CrearJuegoPasaBolas — Crea un nuevo juego de Pasa Bolas
func CrearJuegoPasaBolas(c *gin.Context) {
	// Validar que el cuerpo de la solicitud no esté vacío
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la solicitud vacío"})
		return
	}

	// Crear un slice para almacenar los jugadores
	var jugadores []models.Jugador

	// Obtener los jugadores del cuerpo de la solicitud
	if err := c.BindJSON(&jugadores); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Asegurarse de que haya exactamente 4 jugadores
	if len(jugadores) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe haber exactamente 4 jugadores"})
		return
	}

	// Inicializar los jugadores con 10 bolas cada uno y no eliminados
	jugadoresPasaBolas := make([]models.JugadorPasaBolas, len(jugadores))
	for i, jugador := range jugadores {
		jugadoresPasaBolas[i] = models.JugadorPasaBolas{
			Jugador:   jugador,
			Bolas:     10, // Cada jugador comienza con 10 bolas
			Eliminado: false,
		}
	}

	// Crear un nuevo juego
	juego := models.PasaBolas{
		ID:           fmt.Sprintf("%d", time.Now().UnixNano()), // Generar un ID único con timestamp
		TipoJuego:    "pasa_bolas",
		Jugadores:    jugadoresPasaBolas,
		Estado:       "En Progreso",
		Ciclos:       0,
		Temporizador: 60, // Ciclo de 60 segundos (predeterminado)
	}

	// Guardar el juego en memoria
	juegosMutexPasaBolas.Lock()
	juegosActivosPasaBolas[juego.ID] = juego
	juegosMutexPasaBolas.Unlock()

	// Responder con el juego creado
	c.JSON(http.StatusCreated, gin.H{
		"message": "Juego creado",
		"juego":   juego,
	})
}

// ObtenerJuegoPasaBolas — Devuelve el estado del juego actual
func ObtenerJuegoPasaBolas(c *gin.Context) {
	id := c.Param("id")

	// Acceder al juego de manera segura
	juegosMutexPasaBolas.RLock()
	juego, existe := juegosActivosPasaBolas[id]
	juegosMutexPasaBolas.RUnlock()

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"juego": juego})
}

// LanzarBola — Lanza una bola de un jugador a otro
func LanzarBola(c *gin.Context) {
	id := c.Param("id")
	var movimiento struct {
		DesdeID uint `json:"desde_id"` // ID del jugador que lanza la bola
		HaciaID uint `json:"hacia_id"` // ID del jugador que recibe la bola
	}

	if err := c.BindJSON(&movimiento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Acceder al juego para modificarlo
	juegosMutexPasaBolas.Lock()
	juego, existe := juegosActivosPasaBolas[id]
	if !existe {
		juegosMutexPasaBolas.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Verificar que los jugadores existan y no estén eliminados
	var desdeJugador, haciaJugador *models.JugadorPasaBolas
	for i := range juego.Jugadores {
		if juego.Jugadores[i].ID == movimiento.DesdeID && !juego.Jugadores[i].Eliminado {
			desdeJugador = &juego.Jugadores[i]
		}
		if juego.Jugadores[i].ID == movimiento.HaciaID && !juego.Jugadores[i].Eliminado {
			haciaJugador = &juego.Jugadores[i]
		}
	}

	if desdeJugador == nil || haciaJugador == nil {
		juegosMutexPasaBolas.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jugadores no válidos o eliminados"})
		return
	}

	// Verificar que el jugador que lanza tenga bolas disponibles
	if desdeJugador.Bolas <= 0 {
		juegosMutexPasaBolas.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tienes bolas para lanzar"})
		return
	}

	// Realizar el movimiento
	desdeJugador.Bolas--
	haciaJugador.Bolas++

	// Verificar si el jugador que lanza se quedó sin bolas y eliminarlo
	if desdeJugador.Bolas == 0 {
		desdeJugador.Eliminado = true
	}

	// Verificar si quedan jugadores activos
	activos := 0
	for _, jugador := range juego.Jugadores {
		if !jugador.Eliminado {
			activos++
		}
	}

	// Si queda un solo jugador activo, el juego termina
	if activos == 1 {
		juego.Estado = "Terminado"
	}

	// Actualizar el juego
	juegosActivosPasaBolas[id] = juego
	juegosMutexPasaBolas.Unlock()

	// Responder con el estado actualizado
	c.JSON(http.StatusOK, gin.H{
		"message": "Bola lanzada",
		"juego":   juego,
	})
}

// TerminarJuegoPasaBolas — Termina un juego activo y lo elimina
func TerminarJuegoPasaBolas(c *gin.Context) {
	id := c.Param("id")

	juegosMutexPasaBolas.Lock()
	_, existe := juegosActivosPasaBolas[id]
	if !existe {
		juegosMutexPasaBolas.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}
	delete(juegosActivosPasaBolas, id)
	juegosMutexPasaBolas.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Juego terminado y eliminado"})
}

// ReiniciarBolasPasaBolas — Reinicia las bolas para todos los jugadores
func ReiniciarBolasPasaBolas(c *gin.Context) {
	id := c.Param("id")

	juegosMutexPasaBolas.Lock()
	juego, existe := juegosActivosPasaBolas[id]
	if !existe {
		juegosMutexPasaBolas.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Reiniciar bolas y eliminar a todos los jugadores
	for i := range juego.Jugadores {
		juego.Jugadores[i].Bolas = 10
		juego.Jugadores[i].Eliminado = false
	}

	juego.Estado = "En Progreso"
	juegosActivosPasaBolas[id] = juego
	juegosMutexPasaBolas.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Bolas reiniciadas", "juego": juego})
}
