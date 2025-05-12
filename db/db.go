package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"juego/models"
	"log"
)

var DB *gorm.DB
var err error

func Connect() {
	// Establecer la conexión con la base de datos SQLite
	DB, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		log.Fatal("Error al conectar con la base de datos: ", err)
	}

	// Migrar los modelos
	DB.AutoMigrate(&models.Jugador{})
}

func Close() {
	DB.Close()
}
