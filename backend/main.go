package main

import (
	"fmt"
	"log"
	"net/http"
	"olhar-biblico/config"
	"olhar-biblico/router"
)

func main() {
	config.Load()

	r := router.Config()

	if config.Production {
		fmt.Printf("Production server running on port %d.\n", config.Port)
		if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Port), config.CertFile, config.KeyFile, r); err != nil {
			log.Fatal("Erro:", err.Error())
		}
	}

	fmt.Printf("Development server running on port %d.\n", config.Port)
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", config.Port), r); err != nil {
		log.Fatalf("Error on start server: %s", err.Error())
	}
}
