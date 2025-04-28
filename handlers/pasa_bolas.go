package handlers

import (
	"fmt"
	"juego/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var juegosActivosPasaBolas = make(map[string]models.PasaBolas)

// CrearJuegoPasaBolas crea un nuevo juego de Pasa Bolas
func CrearJuegoPasaBolas(c *gin.Context) {
	// Validar que el cuerpo de la solicitud no esté vacío
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la solicitud vacío"})
		return
	}
	// Crear una estructura para almacenar los jugadores y el temporizador
	var req struct {
		Jugadores    []models.JugadorPasaBolas `json:"jugadores"`
		Temporizador int                       `json:"temporizador"` // Permitir que se pase el tiempo de ciclo
	}

	// Obtener los jugadores y el temporizador del cuerpo de la solicitud
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Validación de número de jugadores
	if len(req.Jugadores) < 2 || len(req.Jugadores) > 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe haber al menos 2 jugadores y un máximo de 6 jugadores"})
		return
	}

	// Asignar bolas a cada jugador
	for i := range req.Jugadores {
		req.Jugadores[i].Bolas = 10  // Inicializar el número de bolas a 10
		req.Jugadores[i].Eliminado = false // Inicializar el estado de eliminación a falso
	}

	// Crear el juego
	juego := models.PasaBolas{
		ID:           fmt.Sprintf("%d", time.Now().Unix()),
		TipoJuego:    "pasa_bolas",
		Jugadores:    req.Jugadores, // Los jugadores asignados desde la solicitud
		Estado:       "En Progreso",
		Ciclos:       0,              // Número de ciclos de eliminación
		Temporizador: req.Temporizador, // El tiempo que se pasa en la solicitud
	}

	// Calcular el número de ciclos basados en la cantidad de jugadores
	juego.Ciclos = len(req.Jugadores) - 1

	// Guardar el juego en memoria (mapa de juegos activos)
	juegosActivosPasaBolas[juego.ID] = juego

	// Responder con el juego creado
	c.JSON(http.StatusCreated, gin.H{
		"message": "Juego creado",
		"juego":   juego,
	})
}


// ObtenerJuegoPasaBolas obtiene el estado de un juego de pasa bolas por su ID
func ObtenerJuegoPasaBolas(c *gin.Context) {
	// Validar que el ID del juego no esté vacío
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del juego vacío"})
		return
	}
	
	id := c.Param("id")

	// Acceder al juego con mutex de lectura
	juego, existe := juegosActivosPasaBolas[id]
	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Responder con el estado del juego
	c.JSON(http.StatusOK, gin.H{
		"juego": juego,
	})
}

// EliminarJugadorPasaBolas elimina un jugador de un juego de pasa bolas
func EliminarJugadorPasaBolas(c *gin.Context) {
	// Validar que el ID del juego no esté vacío
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del juego vacío"})
		return
	}
	id := c.Param("id")

	// Validar que el cuerpo de la solicitud no esté vacío
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la solicitud vacío"})
		return
	}
	juego, existe := juegosActivosPasaBolas[id]
	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Eliminar al jugador con más bolas
	maxBolas := -1
	// Buscar al jugador con más bolas que no esté eliminado
	// y que no sea el último jugador restante
	var jugadorAEliminar models.JugadorPasaBolas
	for _, jugador := range juego.Jugadores {
		if jugador.Bolas > maxBolas && !jugador.Eliminado {
			maxBolas = jugador.Bolas
			jugadorAEliminar = jugador
		}
	}

	// Marcar como eliminado
	for i := range juego.Jugadores {
		if juego.Jugadores[i].ID == jugadorAEliminar.ID {
			juego.Jugadores[i].Eliminado = true
			break
		}
	}

	// Actualizar el estado del juego
	c.JSON(http.StatusOK, gin.H{
		"message":       "Jugador eliminado",
		"jugador_eliminado": jugadorAEliminar,
		"juego":         juego,
	})
}

// TerminarJuegoPasaBolas termina el juego de pasa bolas
func TerminarJuegoPasaBolas(c *gin.Context) {
	id := c.Param("id")

	_, existe := juegosActivosPasaBolas[id]
	if !existe {
		c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
		return
	}

	// Eliminar el juego de la memoria
	// (en este caso, simplemente lo eliminamos del mapa de juegos activos)
	delete(juegosActivosPasaBolas, id)

	// Responder con el mensaje de juego terminado y eliminado
	c.JSON(http.StatusOK, gin.H{
		"message": "Juego terminado y eliminado",
	})
}
// PasarBola permite que un jugador pase una bola a cualquier otro jugador
func PasarBola(c *gin.Context) {
    idJuego := c.Param("id")
    var req struct {
        JugadorEmisor  string `json:"jugador_emisor"`  // ID del jugador que pasa la bola
        JugadorDestino string `json:"jugador_destino"` // ID del jugador que recibe la bola
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    juego, existe := juegosActivosPasaBolas[idJuego]
    if !existe {
        c.JSON(http.StatusNotFound, gin.H{"error": "Juego no encontrado"})
        return
    }

    // Buscar al jugador emisor
    var emisor *models.JugadorPasaBolas
    for i := range juego.Jugadores {
        if juego.Jugadores[i].ID == req.JugadorEmisor && !juego.Jugadores[i].Eliminado {
            emisor = &juego.Jugadores[i]
            break
        }
    }
    if emisor == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Jugador emisor no encontrado o eliminado"})
        return
    }

    // Buscar al jugador destino
    var destino *models.JugadorPasaBolas
    for i := range juego.Jugadores {
        if juego.Jugadores[i].ID == req.JugadorDestino && !juego.Jugadores[i].Eliminado {
            destino = &juego.Jugadores[i]
            break
        }
    }
	// Verificar si el jugador destino existe y no está eliminado
    if destino == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Jugador destino no encontrado o eliminado"})
        return
    }

    // Evitar que el jugador se pase la bola a sí mismo
    if req.JugadorEmisor == req.JugadorDestino {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No puedes pasarte la bola a ti mismo"})
        return
    }

    // Verificar que el emisor tenga bolas para pasar
    if emisor.Bolas <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "El jugador emisor no tiene bolas para pasar"})
        return
    }

    // Actualizar las bolas: el emisor pierde una y el destino gana una
    emisor.Bolas -= 1
    destino.Bolas += 1

    juegosActivosPasaBolas[idJuego] = juego

    c.JSON(http.StatusOK, gin.H{
        "message":         "Bola pasada exitosamente",
        "jugador_emisor":  emisor,
        "jugador_destino": destino,
        "juego":           juego,
    })
}
