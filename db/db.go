package db

import (
	"fmt"
	"log"
	"juego/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	username := "root"
	password := ""
	host := "127.0.0.1"
	port := "3306"
	dbname := "juego"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Error al conectar con la base de datos:", err)
	}

	err = DB.AutoMigrate(&models.Jugador{})
	if err != nil {
		log.Fatal("❌ Error al migrar modelos:", err)
	}

	log.Println("✅ Conexión exitosa a la base de datos.")
}
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}