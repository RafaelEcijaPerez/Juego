package handlers

import (
	"fmt"
	"juego/models"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var juegosActivosPasaBolas = make(map[string]models.PasaBolas) // Juegos activos
var juegosMutexPasaBolas sync.Mutex

// CrearJuegoPasaBolas — Crea un nuevo juego de Pasa Bolas
func CrearJuegoPasaBolas(c *gin.Context) {
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la solicitud vacío"})
		return
	}

	// Obtener los jugadores del cuerpo de la solicitud
	var jugadores []models.Jugador
	if err := c.BindJSON(&jugadores); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Asegurarse de que haya exactamente 4 jugadores
	if len(jugadores) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe haber exactamente 4 jugadores"})
		return
	}

	// Inicializar las bolas para cada jugador
	var jugadoresPasaBolas []models.JugadorPasaBolas
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i, jugador := range jugadores {
		var bolas []models.Bola
		for j := 0; j < 10; j++ {
			bolas = append(bolas, models.Bola{
				X:        150 + random.Float64()*100,
				Y:        150 + random.Float64()*100,
				VX:       0,
				VY:       0,
				Color:    fmt.Sprintf("#%06X", random.Intn(0xFFFFFF)),
				PlayerID: jugador.ID,
			})
		}

		jugadoresPasaBolas = append(jugadoresPasaBolas, models.JugadorPasaBolas{
			Jugador:   jugador,
			Bolas:     bolas,
			Eliminado: false,
			Posicion:  []string{"arriba", "derecha", "abajo", "izquierda"}[i], // Asignamos posición
		})
	}

	// Crear un nuevo juego
	juego := models.PasaBolas{
		ID:           fmt.Sprintf("%d", time.Now().UnixNano()),
		TipoJuego:    "pasa_bolas",
		Jugadores:    jugadoresPasaBolas,
		Estado:       "En Progreso",
		Ciclos:       0,
		Temporizador: 60, // Ciclo de 60 segundos (predeterminado)
	}

	// Guardar el juego en memoria
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
	juego, existe := juegosActivosPasaBolas[id]

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
	juego, existe := juegosActivosPasaBolas[id]
	if !existe {
		juegosMutexPasaBolas.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	var desdeJugador, haciaJugador *models.JugadorPasaBolas
	for i := range juego.Jugadores {
		if juego.Jugadores[i].Jugador.ID == movimiento.DesdeID && !juego.Jugadores[i].Eliminado {
			desdeJugador = &juego.Jugadores[i]
		}
		if juego.Jugadores[i].Jugador.ID == movimiento.HaciaID && !juego.Jugadores[i].Eliminado {
			haciaJugador = &juego.Jugadores[i]
		}
	}

	if desdeJugador == nil || haciaJugador == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jugadores no válidos o eliminados"})
		return
	}

	// Lógica para lanzar la bola
	if len(desdeJugador.Bolas) > 0 {
		// Asumimos que lanzamos la primera bola (puedes modificar esto)
		bola := &desdeJugador.Bolas[0]
		// Mover bola del jugador que lanza al que recibe
		haciaJugador.Bolas = append(haciaJugador.Bolas, *bola)
		// Eliminar la bola del jugador que lanza
		desdeJugador.Bolas = desdeJugador.Bolas[1:]
	}

	// Actualizar estado del juego
	juegosActivosPasaBolas[id] = juego

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
		var bolas []models.Bola
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		for j := 0; j < 10; j++ {
			bolas = append(bolas, models.Bola{
				X:        150 + random.Float64()*100,
				Y:        150 + random.Float64()*100,
				VX:       0,
				VY:       0,
				Color:    fmt.Sprintf("#%06X", random.Intn(0xFFFFFF)),
				PlayerID: juego.Jugadores[i].Jugador.ID,
			})
		}
		juego.Jugadores[i].Bolas = bolas
		juego.Jugadores[i].Eliminado = false
	}

	juego.Estado = "En Progreso"
	juegosActivosPasaBolas[id] = juego
	juegosMutexPasaBolas.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Bolas reiniciadas", "juego": juego})
}
