package main

import (
	"fmt"
	"log"
	"juego/routes"
)

func main() {
	// Usar el enrutador de Gin
	r := routes.SetupRouter()

	fmt.Println("Servidor ejecutándose en el puerto 8080...")
	err := r.Run(":8080") // Esto ya maneja ListenAndServe internamente
	if err != nil {
		log.Fatal(err)
	}
}

