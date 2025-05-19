package controllers

import (
	"fmt"
	"juego/models"
	"juego/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JugadorController struct {
	JugadorService *services.JugadorService
}

func NewJugadorController(jugadorService *services.JugadorService) *JugadorController {
	return &JugadorController{JugadorService: jugadorService}
}

func (controller *JugadorController) Register(c *gin.Context) {
	var jugador models.Jugador

	// Leer JSON
	if err := c.ShouldBindJSON(&jugador); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no válidos"})
		return
	}

	// 🔍 DEBUG - imprimir los campos recibidos
	fmt.Printf("DEBUG: Name=%s, Email=%s, Password=%s\n", jugador.Name, jugador.Email, jugador.Password)

	// Validación extra opcional
	if jugador.Name == "" || jugador.Email == "" || jugador.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan campos obligatorios"})
		return
	}

	// Registrar
	createdJugador, err := controller.JugadorService.RegisterJugador(&jugador)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"jugador": createdJugador})
}



func (controller *JugadorController) Login(c *gin.Context) {
	var loginInput struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no válidos"})
		return
	}

	jugador, err := controller.JugadorService.LoginJugador(loginInput.Email, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jugador": jugador})
}
