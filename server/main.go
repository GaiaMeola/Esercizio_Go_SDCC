package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
)

type MyService struct {
	mu sync.Mutex
}

// Servizio Stateful: Contatore su file condiviso
func (s *MyService) DoWork(args int, reply *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Leggi dal file "condiviso"
	data, _ := os.ReadFile("../state/counter.txt")
	count, _ := strconv.Atoi(strings.TrimSpace(string(data)))

	count++
	
	// Scrivi il nuovo valore
	os.WriteFile("../state/counter.txt", []byte(strconv.Itoa(count)), 0644)
	
	*reply = count
	return nil
}

func main() {
	port := os.Args[1] // Passiamo la porta come argomento: es. 8001
	weight, _ := strconv.Atoi(os.Args[2])

	rpc.Register(new(MyService))
	l, _ := net.Listen("tcp", ":"+port)

	// Registrazione automatica
	client, _ := rpc.Dial("tcp", "localhost:5000")
	var ok bool
	client.Call("Registry.Register", struct{Addr string; Weight int}{"localhost:"+port, weight}, &ok)

	log.Printf("Server attivo sulla porta %s...", port)
	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}