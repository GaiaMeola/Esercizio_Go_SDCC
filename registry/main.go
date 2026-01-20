package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"
)

type ServiceInfo struct {
	Addr   string
	Weight int
}

type Registry struct {
	mu      sync.Mutex
	servers map[string]ServiceInfo
}

func (r *Registry) Register(s ServiceInfo, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.servers == nil { r.servers = make(map[string]ServiceInfo) }
	r.servers[s.Addr] = s
	log.Printf("Registrato server: %s con peso %d", s.Addr, s.Weight)
	*reply = true
	return nil
}

func (r *Registry) GetServers(args struct{}, reply *[]ServiceInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, info := range r.servers {
		*reply = append(*reply, info)
	}
	return nil
}

func main() {
	reg := &Registry{servers: make(map[string]ServiceInfo)}
	rpc.Register(reg)
	l, _ := net.Listen("tcp", ":5000")
	log.Println("Registry RPC attivo sulla porta 5000...")
	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}