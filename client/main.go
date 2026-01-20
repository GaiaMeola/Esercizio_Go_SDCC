package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"service-registry-go/common"
	"time"
)

type Cache struct {
	Servers    []common.ServiceInfo
	LastUpdate time.Time
	TTL        time.Duration
}

var localCache Cache
var globalConfig common.Config

func getServers() []common.ServiceInfo {
	// Usa il TTL caricato dal JSON
	if time.Since(localCache.LastUpdate) < localCache.TTL && len(localCache.Servers) > 0 {
		return localCache.Servers
	}

	// Usa l'indirizzo del Registry caricato dal JSON
	regClient, err := rpc.Dial("tcp", globalConfig.RegistryAddr)
	if err != nil {
		log.Println("Errore: Registry non raggiungibile a", globalConfig.RegistryAddr)
		return localCache.Servers
	}
	defer regClient.Close()

	var serverList []common.ServiceInfo
	err = regClient.Call("Registry.GetServers", struct{}{}, &serverList)
	if err != nil {
		return localCache.Servers
	}

	localCache.Servers = serverList
	localCache.LastUpdate = time.Now()
	log.Printf("Cache aggiornata: %d server trovati", len(serverList))
	return serverList
}

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
	// 1. Carica la configurazione dal file JSON
	var err error
	globalConfig, err = common.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Impossibile caricare config.json: ", err)
	}

	// 2. Inizializza la Cache con il TTL preso dal file
	localCache = Cache{
		TTL: time.Duration(globalConfig.ClientSettings.CacheTTL) * time.Second,
	}

	rand.Seed(time.Now().UnixNano())
	log.Printf("Client avviato. Registry: %s, TTL Cache: %v", globalConfig.RegistryAddr, localCache.TTL)

	for {
		servers := getServers()
		if len(servers) == 0 {
			log.Println("Nessun server disponibile...")
			time.Sleep(5 * time.Second)
			continue
		}

		target := selectWeightedServer(servers)
		fmt.Printf("\n--- Richiesta a: %s (Peso: %d) ---\n", target.Addr, target.Weight)

		client, err := rpc.Dial("tcp", target.Addr)
		if err != nil {
			log.Printf("Server %s offline, svuoto cache...", target.Addr)
			localCache.Servers = nil
			continue
		}

		// Test Stateless
		var replyStateless common.Reply
		client.Call("MyService.Add", common.ArgsStateless{A: 10, B: 20}, &replyStateless)
		fmt.Printf("[Stateless] 10+20 = %d\n", replyStateless.Result)

		// Test Stateful
		var replyStateful common.Reply
		client.Call("MyService.Increment", common.ArgsStateful{Value: 1}, &replyStateful)
		fmt.Printf("[Stateful] Contatore Globale: %d\n", replyStateful.Result)

		client.Close()
		time.Sleep(3 * time.Second)
	}
}