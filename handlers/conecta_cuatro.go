package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"juego/models"
	"net/http"
	"sync"
	"time"
)

var (
	// Mapa para almacenar juegos activos y su sincronización con mutex
	juegosActivosConecta = make(map[string]models.ConectaCuatro)
	juegosMutexConecta   sync.RWMutex
)

// CrearJuegoConecta Crea un nuevo juego
func CrearJuegoConecta(c *gin.Context) {
	var jugadores []models.Jugador

	if c.BindJSON(&jugadores) != nil || len(jugadores) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe proporcionar exactamente 2 jugadores"})
		return
	}

	juego := models.ConectaCuatro{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
		TipoJuego:   "Conecta_Cuatro",
		Jugadores:   jugadores,
		Tablero:     [6][7]string{},
		Estado:      "En Progreso",
		Turno:       0,
		CreadoEn:    time.Now(),
		Actualizado: time.Now(),
	}

	juegosMutexConecta.Lock()
	juegosActivosConecta[juego.ID] = juego
	juegosMutexConecta.Unlock()

	c.JSON(http.StatusCreated, gin.H{"message": "Juego creado", "juego": juego})
}

// ObtenerJuegoConecta Obtiene el estado de un juego por su ID
func ObtenerJuegoConecta(c *gin.Context) {
	id := c.Param("id")

	juegosMutexConecta.RLock()
	juego, existe := juegosActivosConecta[id]
	juegosMutexConecta.RUnlock()

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"juego": juego})
}

// TerminarJuegoConecta Termina un juego y lo elimina de la memoria
func TerminarJuegoConecta(c *gin.Context) {
	id := c.Param("id")

	juegosMutexConecta.Lock()
	_, existe := juegosActivosConecta[id]
	if !existe {
		juegosMutexConecta.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}
	delete(juegosActivosConecta, id)
	juegosMutexConecta.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Juego terminado y eliminado"})
}

// HacerMovimientoConecta Maneja el movimiento de un jugador
func HacerMovimientoConecta(c *gin.Context) {
	id := c.Param("id")
	var movimiento struct {
		Columna int `json:"columna"`
	}

	if c.BindJSON(&movimiento) != nil || movimiento.Columna < 0 || movimiento.Columna >= 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movimiento inválido"})
		return
	}

	juegosMutexConecta.Lock()
	juego, existe := juegosActivosConecta[id]
	if !existe {
		juegosMutexConecta.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Colocar la ficha
	columnaLlena := true
	var filaInsertada int
	for fila := 5; fila >= 0; fila-- {
		if juego.Tablero[fila][movimiento.Columna] == "" {
			filaInsertada = fila
			columnaLlena = false
			if juego.Turno == 0 {
				juego.Tablero[fila][movimiento.Columna] = "X"
			} else {
				juego.Tablero[fila][movimiento.Columna] = "O"
			}
			break
		}
	}

	if columnaLlena {
		juegosMutexConecta.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{"error": "La columna está llena"})
		return
	}

	// Revisar victoria
	ficha := "X"
	if juego.Turno == 1 {
		ficha = "O"
	}
	if revisarVictoria(juego.Tablero, filaInsertada, movimiento.Columna, ficha) {
		juego.Estado = fmt.Sprintf("¡Jugador %d ha ganado!", juego.Turno+1)
		juegosActivosConecta[id] = juego
		juegosMutexConecta.Unlock()
		c.JSON(http.StatusOK, gin.H{"message": juego.Estado, "juego": juego})
		return
	}

	// Revisar empate
	if tableroLleno(juego.Tablero) {
		juego.Estado = "Empate"
		juegosActivosConecta[id] = juego
		juegosMutexConecta.Unlock()
		c.JSON(http.StatusOK, gin.H{"message": "El juego terminó en empate", "juego": juego})
		return
	}

	// Cambiar turno
	juego.Turno = (juego.Turno + 1) % 2
	juego.Actualizado = time.Now()
	juegosActivosConecta[id] = juego
	juegosMutexConecta.Unlock()

	c.JSON(http.StatusOK, gin.H{"juego": juego})
}

// revisarVictoria Verifica si el movimiento actual genera una línea ganadora
func revisarVictoria(tablero [6][7]string, fila, columna int, ficha string) bool {
	direcciones := [][]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}

	for _, dir := range direcciones {
		conteo := 1
		for paso := 1; paso <= 3; paso++ {
			nuevaFila := fila + paso*dir[0]
			nuevaColumna := columna + paso*dir[1]
			if nuevaFila >= 0 && nuevaFila < 6 && nuevaColumna >= 0 && nuevaColumna < 7 && tablero[nuevaFila][nuevaColumna] == ficha {
				conteo++
			} else {
				break
			}
		}
		for paso := 1; paso <= 3; paso++ {
			nuevaFila := fila - paso*dir[0]
			nuevaColumna := columna - paso*dir[1]
			if nuevaFila >= 0 && nuevaFila < 6 && nuevaColumna >= 0 && nuevaColumna < 7 && tablero[nuevaFila][nuevaColumna] == ficha {
				conteo++
			} else {
				break
			}
		}
		if conteo >= 4 {
			return true
		}
	}
	return false
}

// tableroLleno Verifica si el tablero está completamente lleno
func tableroLleno(tablero [6][7]string) bool {
	for _, fila := range tablero {
		for _, celda := range fila {
			if celda == "" {
				return false
			}
		}
	}
	return true
}
