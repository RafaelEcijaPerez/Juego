package models
import ()

type Jugador struct {
	ID   string `json:"id" bson:"id"`   // ID único del jugador
	Nombre string `json:"nombre" bson:"nombre"` // Nombre del jugador
	Contrasena string `json:"contrasena" bson:"contrasena"` // Contraseña del jugador
}