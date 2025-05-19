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
	juegosActivos = make(map[string]models.CuatroEnRaya)
	juegosMutex   sync.RWMutex
)

// CrearJuego — Crea un nuevo juego de Cuatro en Raya
func CrearJuego(c *gin.Context) {
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
		TipoJuego:   "4_en_raya",
		Jugadores:   jugadores,
		Tablero:     [4][4]string{},
		Estado:      "En Progreso",
		Turno:       0, // Comienza el jugador 0
		CreadoEn:    time.Now(),
		Actualizado: time.Now(),
	}

	// Guardar el juego en memoria (con mutex para evitar condiciones de carrera)
	juegosMutex.Lock()
	juegosActivos[juego.ID] = juego
	juegosMutex.Unlock()

	// Responder con el juego creado
	c.JSON(http.StatusCreated, gin.H{
		"message": "Juego creado",
		"juego":   juego,
	})
}

// ObtenerJuego — Obtiene el estado de un juego por su ID
func ObtenerJuego(c *gin.Context) {
	id := c.Param("id")

	// Acceder al juego con mutex de lectura
	juegosMutex.RLock()
	juego, existe := juegosActivos[id]
	juegosMutex.RUnlock()

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"juego": juego})
}

// TerminarJuego — Termina un juego y lo elimina de la memoria
func TerminarJuego(c *gin.Context) {
	id := c.Param("id")

	juegosMutex.Lock()
	_, existe := juegosActivos[id]
	if !existe {
		juegosMutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}
	delete(juegosActivos, id)
	juegosMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Juego terminado y eliminado"})
}

// HacerMovimiento — Maneja el movimiento de un jugador en el juego (colocar o mover ficha)
func HacerMovimiento(c *gin.Context) {
	id := c.Param("id")
	var movimiento struct {
		OrigenX  int `json:"origen_x,omitempty"` // Posición X de la ficha que quiere mover (opcional)
		OrigenY  int `json:"origen_y,omitempty"` // Posición Y de la ficha que quiere mover (opcional)
		DestinoX int `json:"destino_x"`          // Posición X del destino
		DestinoY int `json:"destino_y"`          // Posición Y del destino
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
	juegosMutex.Lock()
	juego, existe := juegosActivos[id]
	if !existe {
		juegosMutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Identificar ficha del jugador actual
	ficha := "X"
	if juego.Turno == 1 {
		ficha = "O"
	}

	// Contar cuántas fichas tiene el jugador en el tablero
	contadorFichas := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if juego.Tablero[i][j] == ficha {
				contadorFichas++
			}
		}
	}

	// Si el jugador tiene menos de 4 fichas, está en la fase de colocación
	if contadorFichas < 4 {
		// Verificar que el destino esté vacío
		if juego.Tablero[movimiento.DestinoX][movimiento.DestinoY] != "" {
			juegosMutex.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "La celda de destino ya está ocupada"})
			return
		}
		// Colocar la ficha en la posición de destino
		juego.Tablero[movimiento.DestinoX][movimiento.DestinoY] = ficha
	} else {
		// Si ya tiene 4 fichas, estamos en la fase de movimiento
		if movimiento.OrigenX < 0 || movimiento.OrigenX >= 4 || movimiento.OrigenY < 0 || movimiento.OrigenY >= 4 {
			juegosMutex.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Posición de origen fuera del tablero"})
			return
		}
		// Verificar que la posición de origen contiene una ficha del jugador
		if juego.Tablero[movimiento.OrigenX][movimiento.OrigenY] != ficha {
			juegosMutex.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "La celda de origen no contiene tu ficha"})
			return
		}
		// Verificar que el destino esté vacío
		if juego.Tablero[movimiento.DestinoX][movimiento.DestinoY] != "" {
			juegosMutex.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "La celda de destino ya está ocupada"})
			return
		}
		// Mover la ficha de origen a destino
		juego.Tablero[movimiento.OrigenX][movimiento.OrigenY] = ""
		juego.Tablero[movimiento.DestinoX][movimiento.DestinoY] = ficha
	}

	// Alternar el turno
	juego.Turno = 1 - juego.Turno

	// Verificar si hay un ganador
	if verificarVictoria(juego.Tablero, ficha) {
		juego.Estado = "Terminado"
		juego.Ganador = &juego.Jugadores[juego.Turno]
		juegosActivos[id] = juego
		juegosMutex.Unlock()
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("¡Jugador %d ha ganado!", juego.Turno+1), "juego": juego})
		return
	}

	// Actualizar el estado del tablero
	juegosActivos[id] = juego
	juegosMutex.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Movimiento realizado", "juego": juego})
}

// verificarVictoria — Revisa si un jugador ha ganado
func verificarVictoria(tablero [4][4]string, ficha string) bool {
	direcciones := [][]int{{0, 1}, {1, 0}, {1, 1}, {-1, 1}} // Horizontal, vertical, diagonales
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
