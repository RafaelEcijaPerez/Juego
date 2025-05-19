package models

import (
	"time"
)

type JugadorConectaCuatro struct {
	Jugador          // Hereda de Jugador
	Ficha     string `json:"ficha"`      // Ficha del jugador (X o O)
	PosicionX int    `json:"posicion_x"` // Posición en el eje X del tablero
	PosicionY int    `json:"posicion_y"` // Posición en el eje Y del tablero
}

type ConectaCuatro struct {
	ID          string       `json:"id"`
	TipoJuego   string       `json:"tipo_juego"`
	Jugadores   []Jugador    `json:"jugadores"`
	Tablero     [6][7]string `json:"tablero"`
	Estado      string       `json:"estado"`
	CreadoEn    time.Time    `json:"creado_en"`
	Actualizado time.Time    `json:"actualizado_en"`
	Turno       int          `json:"turno"`            // 0 para el primer jugador, 1 para el segundo
	Ganador     *Jugador     `json:"winner,omitempty"` // Jugador ganador (si existe)
}
