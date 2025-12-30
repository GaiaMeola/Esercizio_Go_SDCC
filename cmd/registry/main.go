package main

import (
    "log"
    "net/http"

    "service-registry-go/internal/registry"
)

func main() {
    // Creiamo lo store (memoria dei servizi)
    store := registry.NewStore()

    // Creiamo l'handler HTTP
    handler := registry.NewHandler(store)

    // Registriamo gli endpoint HTTP
    http.HandleFunc("/register", handler.Register)
    http.HandleFunc("/deregister/", handler.Deregister)
    http.HandleFunc("/services", handler.List)

    // Log di avvio
    log.Println("Service Registry avviato sulla porta 8500")

    // Avviamo il server HTTP (bloccante)
    log.Fatal(http.ListenAndServe(":8500", nil))
}