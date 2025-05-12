package models

// JugadorPasaBolas Jugador representa a un jugador en el juego Pasa Bolas
type JugadorPasaBolas struct {
	Jugador        // Hereda de Jugador
	Bolas     int  `json:"bolas"`     // Número de bolas que tiene el jugador
	Eliminado bool `json:"eliminado"` // Si el jugador ha sido eliminado
}

// PasaBolas representa el estado del juego Pasa Bolas
type PasaBolas struct {
	ID           string             `json:"id"`           // ID único del juego
	TipoJuego    string             `json:"tipo_juego"`   // Tipo de juego
	Jugadores    []JugadorPasaBolas `json:"jugadores"`    // Lista de jugadores
	Estado       string             `json:"estado"`       // Estado del juego
	Ciclos       int                `json:"ciclos"`       // Número de ciclos
	Temporizador int                `json:"temporizador"` // Tiempo de ciclo en segundos
}
