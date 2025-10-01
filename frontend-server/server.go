package main

import (
    "log"
    "net/http"
)

func main() {
    // Sirve los archivos de la carpeta 'templates'
    http.Handle("/", http.FileServer(http.Dir("../templates")))
    log.Println("Frontend server listening on :8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}