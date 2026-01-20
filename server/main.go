package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"service-registry-go/common"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type MyService struct {
	mu sync.Mutex
}

// Servizio Stateless: Somma
func (s *MyService) Add(args common.ArgsStateless, reply *common.Reply) error {
	reply.Result = args.A + args.B
	log.Printf("Richiesta Stateless: %d + %d", args.A, args.B)
	return nil
}

// Servizio Stateful: Contatore condiviso su file
func (s *MyService) Increment(args common.ArgsStateful, reply *common.Reply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := "../state/counter.txt"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("errore lettura stato: %v", err)
	}

	currentVal, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	newVal := currentVal + args.Value

	err = os.WriteFile(filePath, []byte(strconv.Itoa(newVal)), 0644)
	if err != nil {
		return fmt.Errorf("errore scrittura stato: %v", err)
	}

	reply.Result = newVal
	log.Printf("Richiesta Stateful: Incremento di %d. Nuovo totale: %d", args.Value, newVal)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Uso: go run main.go <porta> <peso>")
	}
	port := os.Args[1]
	weight, _ := strconv.Atoi(os.Args[2])
	addr := "localhost:" + port

	// 1. CARICAMENTO CONFIGURAZIONE
	// Leggiamo il JSON per sapere dove inviare la registrazione
	config, err := common.LoadConfig("../config.json")
	if err != nil {
		log.Fatal("Errore caricamento config.json: ", err)
	}

	// 2. Setup RPC locale
	service := new(MyService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Errore avvio listener: %v", err)
	}
	log.Printf("Server attivo su %s (Registry puntato a: %s)", addr, config.RegistryAddr)

	// 3. Connessione al Registry usando l'indirizzo del JSON
	regClient, err := rpc.Dial("tcp", config.RegistryAddr)
	if err != nil {
		log.Fatal("ERRORE: Impossibile connettersi al Registry. Controlla config.json e riprova.")
	}

	// 4. Registrazione automatica
	var ok bool
	regArgs := common.RegistryArgs{Service: common.ServiceInfo{Addr: addr, Weight: weight}}
	err = regClient.Call("Registry.Register", regArgs, &ok)
	if err != nil {
		log.Printf("Errore registrazione: %v", err)
	} else {
		log.Println("Registrato con successo al Registry")
	}

	// 5. Gestione spegnimento pulito (Deregister)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("\nSpegnimento... Deregistrazione in corso")
		var reply bool
		deregArgs := common.RegistryArgs{Service: common.ServiceInfo{Addr: addr}}
		regClient.Call("Registry.Deregister", deregArgs, &reply)
		regClient.Close()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}