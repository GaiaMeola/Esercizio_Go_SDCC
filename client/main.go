package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type ServiceInfo struct {
	Addr   string
	Weight int
}

var cache []ServiceInfo
var lastUpdate time.Time
var rrIndex int = 0

func getServers() []ServiceInfo {
	// Caching dinamica: se l'ultima lookup è più vecchia di 10 secondi, aggiorna
	if time.Since(lastUpdate) < 10*time.Second && len(cache) > 0 {
		return cache
	}

	reg, _ := rpc.Dial("tcp", "localhost:5000")
	var servers []ServiceInfo
	reg.Call("Registry.GetServers", struct{}{}, &servers)
	cache = servers
	lastUpdate = time.Now()
	return cache
}

func main() {
	for {
		servers := getServers()
		// Algoritmo Round Robin semplice
		target := servers[rrIndex % len(servers)]
		rrIndex++

		client, err := rpc.Dial("tcp", target.Addr)
		if err != nil {
			log.Println("Errore server, invalido la cache...")
			cache = nil // Invalidazione cache su errore
			continue
		}

		var result int
		client.Call("MyService.DoWork", 1, &result)
		fmt.Printf("Risposta dal server %s - Contatore globale: %d\n", target.Addr, result)
		
		time.Sleep(2 * time.Second)
	}
}