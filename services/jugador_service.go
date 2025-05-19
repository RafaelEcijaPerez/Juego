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
	// Verificar si el email ya existe
	var existing models.Jugador
	if err := service.DB.Where("email = ?", jugador.Email).First(&existing).Error; err == nil {
		return nil, errors.New("el correo ya está registrado")
	}

	// Generar hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(jugador.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("no se pudo procesar la contraseña")
	}
	jugador.Password = string(hashedPassword)

	// Crear jugador en la base de datos
	if err := service.DB.Create(jugador).Error; err != nil {
		return nil, err
	}

	// Retornar solo los datos visibles
	return &models.Jugador{
		ID:     jugador.ID,
		Name:   jugador.Name,
		Email:  jugador.Email,
	}, nil
}

func (service *JugadorService) LoginJugador(email, password string) (*models.Jugador, error) {
	// Buscar al jugador por correo
	var Jugador models.Jugador
	if err := service.DB.Where("email = ?", email).First(&Jugador).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("correo o contraseña incorrectos")
		}
		return nil, err
	}

	// Verificar que la contraseña sea correcta
	if err := bcrypt.CompareHashAndPassword([]byte(Jugador.Password), []byte(password)); err != nil {
		return nil, errors.New("correo o contraseña incorrectos")
	}

	// Retornar solo la información pública del jugador
	jugador := &models.Jugador{
		ID:     Jugador.ID,
		Name: Jugador.Name,
		Email:  Jugador.Email,
	}

	return jugador, nil
}
