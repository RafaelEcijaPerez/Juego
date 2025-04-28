package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// db será el objeto de conexión a la base de datos
var DB *gorm.DB

// InitDB inicializa la conexión a MySQL
func InitDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/nombre_basededatos?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos")
	}
	fmt.Println("Conexión a MySQL exitosa")
}
