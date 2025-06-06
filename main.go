package main

import (
	"log"
	"net/http"

	"github.com/miledxz/salesforge/db"
	"github.com/miledxz/salesforge/router"
)

func main() {
	db.Connect()
	r := router.SetupRouter()
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
