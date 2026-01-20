package main

import (
	"log"
	"net"
	"net/rpc"
	"service-registry-go/common"
	"sync"
)

type Registry struct {
	mu      sync.Mutex
	servers map[string]common.ServiceInfo
}

// Register aggiunge un server alla lista
func (r *Registry) Register(args common.RegistryArgs, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.servers[args.Service.Addr] = args.Service
	log.Printf("Server registrato: %s (Peso: %d)", args.Service.Addr, args.Service.Weight)
	*reply = true
	return nil
}

// Deregister rimuove un server
func (r *Registry) Deregister(args common.RegistryArgs, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.servers, args.Service.Addr)
	log.Printf("Server rimosso: %s", args.Service.Addr)
	*reply = true
	return nil
}

// GetServers restituisce tutti i server attivi al client
func (r *Registry) GetServers(args struct{}, reply *[]common.ServiceInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, info := range r.servers {
		*reply = append(*reply, info)
	}
	return nil
}

func main() {
	reg := &Registry{
		servers: make(map[string]common.ServiceInfo),
	}
	rpc.Register(reg)

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Errore avvio Registry:", err)
	}
	log.Println("Service Registry in ascolto sulla porta 5000...")
	
	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}