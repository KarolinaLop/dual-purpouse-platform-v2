package main

import (
	"log"

	"github.com/KarolinaLop/dp/web"
)

// Entry point to the app
func main() {
	s := web.SetupServer()

	log.Println("Server is running on http://localhost:" + web.PORT)
	log.Fatal(s.ListenAndServe())
}
