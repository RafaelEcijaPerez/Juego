package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

type Jugador struct {
	ID   string `json:"id" bson:"id"`   // ID único del jugador
	Nombre string `json:"nombre" bson:"nombre"` // Nombre del jugador
	Email    string `json:"email"`
	Contrasena string `json:"contrasena" bson:"contrasena"` // Contraseña del jugador
}

var DB *gorm.DB
var err error

func Connect() {
	// Establecer la conexión con la base de datos SQLite
	DB, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		log.Fatal("Error al conectar con la base de datos: ", err)
	}

	// Migrar los modelos
	DB.AutoMigrate(&Jugador{})
}

func Close() {
	DB.Close()
}