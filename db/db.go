package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"juego/models"
	"log"
)

var DB *gorm.DB
var err error

func Connect() {
	// Primero intenta conectar con la base de datos
	username := "root"
	password := ""
	host := "127.0.0.1"
	port := "3306"
	dbname := "juego"

	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8&parseTime=True&loc=Local"

	// Intentar establecer la conexión
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al conectar con la base de datos: ", err)
	}

	// Verificar que la conexión es válida
	if DB == nil {
		log.Fatal("La conexión a la base de datos no se ha inicializado correctamente.")
	}

	// Desactiva la pluralización de las tablas de GORM
	DB.SingularTable(true)

	// Migrar los modelos
	DB.AutoMigrate(&models.Jugador{})
}

func Close() {
	if err := DB.Close(); err != nil {
		log.Fatal("Error al cerrar la base de datos: ", err)
	}
}
