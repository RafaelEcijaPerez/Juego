package main

import (
	"fmt"
	"juego/db"
	"juego/routes"
	"log"
)

func main() {
	// Conexión a la base de datos
	db.Connect()
	// cierra la conexión
	defer db.Close()

	// Usar el enrutador de Gin
	r := routes.SetupRouter()

	fmt.Println("Servidor ejecutándose en el puerto 8080...")
	err := r.Run(":8081") // Esto ya maneja ListenAndServe internamente
	if err != nil {
		log.Fatal(err)
	}
}
