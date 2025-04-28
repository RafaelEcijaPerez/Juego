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

// CrearJuego crea un nuevo juego de Cuatro en Raya
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

// ObtenerJuego obtiene el estado de un juego por su ID
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

// TerminarJuego termina un juego y lo elimina de la memoria
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

// HacerMovimiento maneja el movimiento de un jugador en el juego
func HacerMovimiento(c *gin.Context) {
	id := c.Param("id")
	var movimiento struct {
		Columna int `json:"columna"` // Columna en la que se desea colocar la ficha
	}

	if err := c.BindJSON(&movimiento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movimiento inválido"})
		return
	}

	// Verificar que la columna sea válida (0 a 3)
	if movimiento.Columna < 0 || movimiento.Columna >= 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Columna inválida"})
		return
	}

	// Acceder al juego con mutex de escritura
	juegosMutex.Lock()
	juego, existe := juegosActivos[id]
	if !existe {
		juegosMutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Colocar la ficha en la columna indicada (comenzamos desde la última fila)
	columnaLlenada := true
	for fila := 3; fila >= 0; fila-- {
		if juego.Tablero[fila][movimiento.Columna] == "" {
			columnaLlenada = false
			if juego.Turno == 0 {
				juego.Tablero[fila][movimiento.Columna] = "X"
			} else {
				juego.Tablero[fila][movimiento.Columna] = "O"
			}
			break
		}
	}

	// Si la columna está llena, no se puede hacer el movimiento
	// (columnaLlenada se mantiene en true si no se encontró un espacio vacío)
	if columnaLlenada {
		juegosMutex.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Columna llena"})
		return
	}

	// Cambiar el turno (alternar entre 0 y 1)
	juego.Turno = 1 - juego.Turno

	// Verificar si hay un ganador
	if ganador := verificarGanador(juego.Tablero); ganador != "" {
		juego.Estado = "Terminado"
		// Buscar el jugador ganador basado en la ficha ("X" o "O")
		var ganadorJugador *models.Jugador
		if ganador == "X" {
			ganadorJugador = &juego.Jugadores[0]
		} else {
			ganadorJugador = &juego.Jugadores[1]
		}
		juego.Ganador = ganadorJugador

		// Actualizar el juego en memoria
		juego.Actualizado = time.Now()
		juegosActivos[id] = juego
		juegosMutex.Unlock()

		// Responder con el estado del juego y el ganador
		// (eliminamos el juego de la memoria para evitar que se pueda jugar nuevamente)
		c.JSON(http.StatusOK, gin.H{
			"message": "Juego terminado",
			"winner":  ganadorJugador,
		})
		return
	}

	// Si el juego continúa, actualizar la fecha de última actualización
	juego.Actualizado = time.Now()
	juegosActivos[id] = juego
	juegosMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Movimiento realizado",
		"juego":   juego,
	})
}

// verificarGanador revisa filas, columnas y diagonales en el tablero para determinar un ganador
func verificarGanador(tablero [4][4]string) string {
	// Revisar filas
	for fila := 0; fila < 4; fila++ {
		if tablero[fila][0] != "" &&
			tablero[fila][0] == tablero[fila][1] &&
			tablero[fila][0] == tablero[fila][2] &&
			tablero[fila][0] == tablero[fila][3] {
			return tablero[fila][0]
		}
	}

	// Revisar columnas
	for col := 0; col < 4; col++ {
		if tablero[0][col] != "" &&
			tablero[0][col] == tablero[1][col] &&
			tablero[0][col] == tablero[2][col] &&
			tablero[0][col] == tablero[3][col] {
			return tablero[0][col]
		}
	}

	// Revisar diagonal principal
	if tablero[0][0] != "" &&
		tablero[0][0] == tablero[1][1] &&
		tablero[0][0] == tablero[2][2] &&
		tablero[0][0] == tablero[3][3] {
		return tablero[0][0]
	}

	// Revisar diagonal secundaria
	if tablero[0][3] != "" &&
		tablero[0][3] == tablero[1][2] &&
		tablero[0][3] == tablero[2][1] &&
		tablero[0][3] == tablero[3][0] {
		return tablero[0][3]
	}

	return ""
}
