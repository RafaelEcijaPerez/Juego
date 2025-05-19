package models

// JugadorPasaBolas Jugador representa a un jugador en el juego Pasa Bolas
type JugadorPasaBolas struct {
	Jugador   Jugador  `json:"jugador"`
	Bolas     []Bola   `json:"bolas"`
	Eliminado bool     `json:"eliminado"`
	Posicion  string   `json:"posicion"`
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
 type Bola struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	VX        float64 `json:"vx"`
	VY        float64 `json:"vy"`
	Color     string  `json:"color"`     // Puedes usar hex string (ej: "#FF00FF")
	PlayerID  uint    `json:"player_id"` // ID del jugador dueño
}
