package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"service-registry-go/common"
	"time"
)

// Struttura per gestire la Cache locale
type Cache struct {
	Servers    []common.ServiceInfo
	LastUpdate time.Time
	TTL        time.Duration
}

var localCache = Cache{
	TTL: 10 * time.Second, // La cache scade dopo 10 secondi
}

// Funzione per ottenere i server (con Service Discovery e Caching)
func getServers() []common.ServiceInfo {
	if time.Since(localCache.LastUpdate) < localCache.TTL && len(localCache.Servers) > 0 {
		return localCache.Servers
	}

	// Se la cache è scaduta, interroga il Registry
	regClient, err := rpc.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Println("Errore: Registry non raggiungibile.")
		return localCache.Servers
	}
	defer regClient.Close()

	var serverList []common.ServiceInfo
	err = regClient.Call("Registry.GetServers", struct{}{}, &serverList)
	if err != nil {
		log.Println("Errore lookup server:", err)
		return localCache.Servers
	}

	localCache.Servers = serverList
	localCache.LastUpdate = time.Now()
	log.Printf("Cache aggiornata: %d server trovati", len(serverList))
	return serverList
}

// Algoritmo Weighted Load Balancing
func selectWeightedServer(servers []common.ServiceInfo) common.ServiceInfo {
	totalWeight := 0
	for _, s := range servers {
		totalWeight += s.Weight
	}
	if totalWeight == 0 {
		return servers[0]
	}
	r := rand.Intn(totalWeight)
	for _, s := range servers {
		r -= s.Weight
		if r < 0 {
			return s
		}
	}
	return servers[0]
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		servers := getServers()
		if len(servers) == 0 {
			log.Println("Nessun server disponibile. Riprovo tra 5 secondi...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Scegliamo un server usando l'algoritmo Weighted
		target := selectWeightedServer(servers)
		fmt.Printf("\n--- Chiamata verso: %s (Peso: %d) ---\n", target.Addr, target.Weight)

		// Connessione RPC al Server scelto
		client, err := rpc.Dial("tcp", target.Addr)
		if err != nil {
			log.Printf("Server %s offline, invalido la cache...", target.Addr)
			localCache.Servers = nil // Invalida la cache per forzare un nuovo lookup
			continue
		}

		// 1. Test Servizio Stateless (Somma)
		argsStateless := common.ArgsStateless{A: 10, B: 20}
		var replyStateless common.Reply
		err = client.Call("MyService.Add", argsStateless, &replyStateless)
		if err == nil {
			fmt.Printf("[Stateless] Somma 10+20 = %d\n", replyStateless.Result)
		}

		// 2. Test Servizio Stateful (Contatore)
		argsStateful := common.ArgsStateful{Value: 1}
		var replyStateful common.Reply
		err = client.Call("MyService.Increment", argsStateful, &replyStateful)
		if err == nil {
			fmt.Printf("[Stateful] Contatore Globale: %d\n", replyStateful.Result)
		}

		client.Close()
		time.Sleep(3 * time.Second) // Attendi prima della prossima chiamata
	}
}