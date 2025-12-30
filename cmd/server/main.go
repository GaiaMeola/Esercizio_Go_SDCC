package main

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
)

type Service struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Host string `json:"host"`
    Port int    `json:"port"`
}

func main() {
    // Parametri da riga di comando
    name := flag.String("name", "Auth", "Nome logico del servizio")
    id := flag.String("id", "", "ID univoco dell'istanza")
    port := flag.Int("port", 8081, "Porta del servizio")
    flag.Parse()

    if *id == "" {
        *id = fmt.Sprintf("%s-%d", *name, *port)
    }

    service := Service{
        ID:   *id,
        Name: *name,
        Host: "localhost",
        Port: *port,
    }

    // Registrazione al registry
    data, _ := json.Marshal(service)
    resp, err := http.Post("http://localhost:8500/register", "application/json", bytes.NewBuffer(data))
    if err != nil {
        log.Fatal("Errore registrazione:", err)
    }
    resp.Body.Close()
    log.Println("Servizio registrato:", service)

    // Endpoint di esempio
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(fmt.Sprintf("Ciao dal servizio %s!", service.Name)))
    })

    addr := fmt.Sprintf(":%d", service.Port)
    log.Println("Server servizio in ascolto su porta", service.Port)
    log.Fatal(http.ListenAndServe(addr, nil))
}