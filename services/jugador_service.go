package services

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"juego/models"
)

type JugadorService struct {
	DB *gorm.DB
}

func NewJugadorService(db *gorm.DB) *JugadorService {
	return &JugadorService{DB: db}
}

func (service *JugadorService) RegisterJugador(jugador *models.Jugador) (*models.Jugador, error) {
	// Verificar si el correo ya existe
	var existingJugador models.Jugador
	if err := service.DB.Where("email = ?", jugador.Email).First(&existingJugador).Error; err == nil {
		return nil, errors.New("el correo ya está registrado")
	}

	// Encriptar la contraseña antes de guardarla
	hash, err := bcrypt.GenerateFromPassword([]byte(jugador.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	jugador.Password = string(hash)

	// Crear el nuevo jugador
	if err := service.DB.Create(jugador).Error; err != nil {
		return nil, err
	}

	return jugador, nil
}
