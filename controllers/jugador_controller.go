package controllers

import (
	"github.com/gin-gonic/gin"
	"juego/models"
	"juego/services"
	"net/http"
)

type JugadorController struct {
	JugadorService *services.JugadorService
}

func NewJugadorController(jugadorService *services.JugadorService) *JugadorController {
	return &JugadorController{JugadorService: jugadorService}
}

func (controller *JugadorController) Register(c *gin.Context) {
	var jugador models.Jugador

	// Vincular el JSON de la solicitud al modelo de jugador
	if err := c.ShouldBindJSON(&jugador); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no válidos"})
		return
	}

	// Llamar al servicio para registrar al jugador
	createdJugador, err := controller.JugadorService.RegisterJugador(&jugador)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Responder con el jugador creado
	c.JSON(http.StatusCreated, gin.H{"jugador": createdJugador})
}
