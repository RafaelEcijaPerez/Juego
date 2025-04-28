package models

import ()

// Jugador representa a un jugador en el juego Pasa Bolas
type JugadorPasaBolas struct {
	Jugador           // Hereda de Jugador
	Bolas      int    `json:"bolas" bson:"bolas"`       // Número de bolas que tiene el jugador
	Eliminado  bool   `json:"eliminado" bson:"eliminado"` // Si el jugador ha sido eliminado
}

// PasaBolas representa el estado del juego Pasa Bolas
type PasaBolas struct {
	ID           string            `json:"id" bson:"id"`           // ID único del juego
	TipoJuego    string            `json:"tipo_juego" bson:"tipo_juego"` // Tipo de juego
	Jugadores    []JugadorPasaBolas `json:"jugadores" bson:"jugadores"`   // Lista de jugadores
	Estado       string            `json:"estado" bson:"estado"`   // Estado del juego
	Ciclos       int               `json:"ciclos" bson:"ciclos"`   // Número de ciclos
	Temporizador int               `json:"temporizador" bson:"temporizador"` // Tiempo de ciclo en segundos
}