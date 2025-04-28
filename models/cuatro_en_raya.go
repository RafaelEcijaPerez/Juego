package models

import (
	"time"
)
type JugadorCuatroEnRaya struct {
	Jugador    // Hereda de Jugador
	Ficha      string `json:"ficha" bson:"ficha"` // Ficha del jugador (X o O)
	PosicionX  int    `json:"posicion_x" bson:"posicion_x"` // Posición en el eje X del tablero
	PosicionY  int    `json:"posicion_y" bson:"posicion_y"` // Posición en el eje Y del tablero
}
// CuatroEnRaya representa el estado del juego de cuatro en raya
type CuatroEnRaya struct {
	ID           string     `json:"id"`
	TipoJuego    string     `json:"tipo_juego"`
	Jugadores    []Jugador  `json:"jugadores"`
	Tablero      [4][4]string `json:"tablero"`
	Estado       string     `json:"estado"`
	CreadoEn     time.Time  `json:"creado_en"`
	Actualizado  time.Time  `json:"actualizado_en"`
	Turno        int        `json:"turno"` // 0 para el primer jugador, 1 para el segundo
	Ganador       *Jugador   `json:"winner,omitempty"` // Jugador ganador (si existe)
}
