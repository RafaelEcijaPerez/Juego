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

type RegisterInput struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (controller *JugadorController) Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no válidos"})
        return
    }

    jugador := models.Jugador{
        Name:     input.Name,
        Email:    input.Email,
        Password: input.Password, // Aquí sí tenemos la contraseña
    }

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

    fmt.Printf("DEBUG LOGIN: email=%s, password=%s\n", loginInput.Email, loginInput.Password)

    jugador, err := controller.JugadorService.LoginJugador(loginInput.Email, loginInput.Password)
    if err != nil {
        fmt.Printf("DEBUG LOGIN ERROR: %v\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"jugador": jugador})
}
