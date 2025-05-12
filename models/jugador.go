package models

type Jugador struct {
	ID       uint   `json:"id"`
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"contrasena"`
}
