package services

import (
    "errors"
    "fmt"
    "juego/models"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type JugadorService struct {
    DB *gorm.DB
}

func NewJugadorService(db *gorm.DB) *JugadorService {
    return &JugadorService{DB: db}
}

func (service *JugadorService) RegisterJugador(jugador *models.Jugador) (*models.Jugador, error) {
    var existing models.RegisterInput
    if err := service.DB.Where("Email = ?", jugador.Email).First(&existing).Error; err == nil {
        return nil, errors.New("el correo ya está registrado")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(jugador.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("no se pudo procesar la contraseña")
    }
    jugador.Password = string(hashedPassword)

    if err := service.DB.Create(jugador).Error; err != nil {
        return nil, err
    }

    return &models.Jugador{
        ID:    jugador.ID,
        Name:  jugador.Name,
        Email: jugador.Email,
    }, nil
}

func (service *JugadorService) LoginJugador(email, password string) (*models.Jugador, error) {
    var jugador models.Jugador
    err := service.DB.Where("Email = ?", email).First(&jugador).Error
    if err != nil {
        fmt.Printf("DEBUG DB ERROR: %v\n", err)
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("correo o contraseña incorrectos")
        }
        return nil, err
    }

    fmt.Printf("DEBUG BD hash: %s\n", jugador.Password)
    fmt.Printf("DEBUG Password recibido: %s\n", password)

    if err := bcrypt.CompareHashAndPassword([]byte(jugador.Password), []byte(password)); err != nil {
        return nil, errors.New("correo o contraseña incorrectos")
    }

    return &models.Jugador{
        ID:    jugador.ID,
        Name:  jugador.Name,
        Email: jugador.Email,
    }, nil
}
