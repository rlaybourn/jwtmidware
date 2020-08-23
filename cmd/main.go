package main

import (
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rlaybourn/jwtmidware"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// 187ed6e4fcb982bdcb21310e00adcbcb5dd8a3f9
func main() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.Handle("/hello", jwtmidware.ServeProtected(http.HandlerFunc(helloHandler)))
	http.HandleFunc("/login", jwtmidware.Login)
	http.HandleFunc("/logout", jwtmidware.Logout)
	log.Println("Listing for requests at http://localhost:8000/hello")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
