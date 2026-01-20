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

	// Legge lo stato dal file condiviso
	data, err := os.ReadFile("../state/counter.txt")
	if err != nil {
		return fmt.Errorf("errore lettura stato: %v", err)
	}

	currentVal, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	newVal := currentVal + args.Value

	// Scrive il nuovo stato
	err = os.WriteFile("../state/counter.txt", []byte(strconv.Itoa(newVal)), 0644)
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

	// Setup RPC
	service := new(MyService)
	rpc.Register(service)

	listener, _ := net.Listen("tcp", ":"+port)
	log.Printf("Server RPC attivo su %s (Peso: %d)", addr, weight)

	// Registrazione automatica al Registry
	regClient, err := rpc.Dial("tcp", "localhost:5000")
	if err == nil {
		var ok bool
		regArgs := common.RegistryArgs{Service: common.ServiceInfo{Addr: addr, Weight: weight}}
		regClient.Call("Registry.Register", regArgs, &ok)
		log.Println("Registrato con successo al Registry")
	}

	// Gestione Deregistrazione allo spegnimento (CTRL+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Spegnimento in corso...")
		var ok bool
		regArgs := common.RegistryArgs{Service: common.ServiceInfo{Addr: addr}}
		regClient.Call("Registry.Deregister", regArgs, &ok)
		os.Exit(0)
	}()

	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}