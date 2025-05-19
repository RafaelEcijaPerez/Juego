package handlers

import (
	"fmt"
	"juego/models"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Mapa en memoria para almacenar juegos activos y mutex para acceso concurrente
var (
	juegosActivosDesdeBorde = make(map[string]models.CuatroEnRaya)
	juegosMutexDesdeBorde   sync.RWMutex
)

// CrearJuegoDesdeBorde — Crea un nuevo juego de Cuatro en Raya desde el borde
func CrearJuegoDesdeBorde(c *gin.Context) {
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

	// Validación de número de jugadores (exactamente 2)
	if len(jugadores) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe haber exactamente 2 jugadores"})
		return
	}

	// Crear el juego con un tablero vacío y turno inicial en 0
	juego := models.CuatroEnRaya{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()), // Usamos UnixNano para mayor unicidad
		TipoJuego:   "4_en_raya_desde_borde",
		Jugadores:   jugadores,
		Tablero:     [4][4]string{},
		Estado:      "En Progreso",
		Turno:       0, // Comienza el jugador 0
		CreadoEn:    time.Now(),
		Actualizado: time.Now(),
	}

	// Guardar el juego en memoria (con mutex para evitar condiciones de carrera)
	juegosMutexDesdeBorde.Lock()
	juegosActivosDesdeBorde[juego.ID] = juego
	juegosMutexDesdeBorde.Unlock()

	// Responder con el juego creado
	c.JSON(http.StatusCreated, gin.H{
		"message": "Juego creado",
		"juego":   juego,
	})
}

// ObtenerJuegoDesdeBorde — Obtiene el estado de un juego por su ID
func ObtenerJuegoDesdeBorde(c *gin.Context) {
	id := c.Param("id")

	// Acceder al juego con mutex de lectura
	juegosMutexDesdeBorde.RLock()
	juego, existe := juegosActivosDesdeBorde[id]
	juegosMutexDesdeBorde.RUnlock()

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"juego": juego})
}

// TerminarJuegoDesdeBorde — Termina un juego y lo elimina de la memoria
func TerminarJuegoDesdeBorde(c *gin.Context) {
	id := c.Param("id")

	juegosMutexDesdeBorde.Lock()
	_, existe := juegosActivosDesdeBorde[id]
	if !existe {
		juegosMutexDesdeBorde.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}
	delete(juegosActivosDesdeBorde, id)
	juegosMutexDesdeBorde.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Juego terminado y eliminado"})
}

// HacerMovimientoDesdeBorde — Maneja el movimiento (colocación) en el juego
func HacerMovimientoDesdeBorde(c *gin.Context) {
	id := c.Param("id")
	var movimiento struct {
		DestinoX int `json:"destino_x"`
		DestinoY int `json:"destino_y"`
	}

	// Validar entrada del movimiento
	if err := c.BindJSON(&movimiento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movimiento inválido"})
		return
	}

	// Verificar que el destino esté dentro del rango válido del tablero
	if movimiento.DestinoX < 0 || movimiento.DestinoX >= 4 || movimiento.DestinoY < 0 || movimiento.DestinoY >= 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Posición de destino fuera del tablero"})
		return
	}

	// Acceder al juego con el mutex de escritura
	juegosMutexDesdeBorde.Lock()
	juego, existe := juegosActivosDesdeBorde[id]
	if !existe {
		juegosMutexDesdeBorde.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Identificar ficha del jugador actual
	ficha := "X"
	if juego.Turno == 1 {
		ficha = "O"
	}

	// Validar si es parte del anillo exterior o interior permitido
	isOuterRing := movimiento.DestinoX == 0 || movimiento.DestinoX == 3 || movimiento.DestinoY == 0 || movimiento.DestinoY == 3

	// Verificar si el anillo exterior está lleno
	outerRingFull := true
	for i := 0; i < 4; i++ {
		if juego.Tablero[0][i] == "" || juego.Tablero[3][i] == "" || juego.Tablero[i][0] == "" || juego.Tablero[i][3] == "" {
			outerRingFull = false
			break
		}
	}

	// Si el borde exterior no está lleno, solo se pueden colocar fichas allí
	if !outerRingFull && !isOuterRing {
		juegosMutexDesdeBorde.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debes colocar la ficha en el borde exterior primero"})
		return
	}

	// Verificar que el destino esté vacío
	if juego.Tablero[movimiento.DestinoX][movimiento.DestinoY] != "" {
		juegosMutexDesdeBorde.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{"error": "La celda de destino ya está ocupada"})
		return
	}

	// Colocar la ficha en el destino
	juego.Tablero[movimiento.DestinoX][movimiento.DestinoY] = ficha

	// Cambiar el turno (alternar entre 0 y 1)
	juego.Turno = 1 - juego.Turno

	// Verificar si hay un ganador
	if ganador := verificarVictoria(juego.Tablero, ficha); ganador {
		juego.Estado = "Terminado"
		juegosActivosDesdeBorde[id] = juego
		juegosMutexDesdeBorde.Unlock()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("¡Jugador %d ha ganado!", juego.Turno+1),
			"juego":   juego,
		})
		return
	}

	// Actualizar el estado del juego
	juegosActivosDesdeBorde[id] = juego
	juegosMutexDesdeBorde.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Movimiento realizado", "juego": juego})
}

// verificarVictoriaBorde — Revisa si un jugador ha ganado
func verificarVictoriaBorde(tablero [4][4]string, ficha string) bool {
	direcciones := [][]int{{0, 1}, {1, 0}, {1, 1}, {-1, 1}} // Horizontal, vertical, diagonal ascendente y descendente
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if tablero[i][j] != ficha {
				continue
			}
			for _, dir := range direcciones {
				conteo := 1
				for paso := 1; paso < 4; paso++ {
					nuevaX := i + paso*dir[0]
					nuevaY := j + paso*dir[1]
					if nuevaX < 0 || nuevaX >= 4 || nuevaY < 0 || nuevaY >= 4 || tablero[nuevaX][nuevaY] != ficha {
						break
					}
					conteo++
				}
				if conteo >= 4 {
					return true
				}
			}
		}
	}
	return false
}
