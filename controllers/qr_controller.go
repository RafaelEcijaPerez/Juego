// controllers/qr_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QRController struct{}

func NewQRController() *QRController {
	return &QRController{}
}

// Genera un código QR (solo devuelve un UUID)
func (qc *QRController) GenerateQR(c *gin.Context) {
	qrCode := uuid.New().String()

	// Aquí podrías guardar el QR en memoria, Redis o una tabla temporal con estado "pendiente"

	c.JSON(http.StatusOK, gin.H{
		"qr_code": qrCode,
	})
}

// Verifica si el QR ha sido usado (simulado)
func (qc *QRController) CheckQRStatus(c *gin.Context) {
	var input struct {
		QRCode string `json:"qr_code"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.QRCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR inválido"})
		return
	}

	// Aquí simularías que aún no está usado
	c.JSON(http.StatusOK, gin.H{"status": "pending"})
}

// Login usando QR
func (qc *QRController) LoginWithQR(c *gin.Context) {
	var input struct {
		QRData    string `json:"qr_data"`
		Timestamp int64  `json:"timestamp"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.QRData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	

	c.JSON(http.StatusOK, gin.H{
		"jugador": gin.H{
			"id":    1,
			"nombre": "JugadorQR",
			"email": "qr@example.com",
		},
	})
}
