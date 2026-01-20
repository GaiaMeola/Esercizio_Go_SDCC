package common

import (
	"encoding/json"
	"os"
)


// Config riflette la struttura del file JSON
type Config struct {
    RegistryAddr   string `json:"registry_addr"`
    ClientSettings struct {
        CacheTTL   int `json:"cache_ttl_seconds"`
        NumClients int `json:"num_clients"`
    } `json:"client_settings"`
    Servers []struct {
        Port   string `json:"port"`
        Weight int    `json:"weight"`
    } `json:"servers"`
}

// LoadConfig legge il file JSON e lo trasforma nella struttura Config
func LoadConfig(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

type ServiceInfo struct {
	Addr   string 
	Weight int    
}

type ArgsStateless struct {
	A, B int
}

type ArgsStateful struct {
	Value int 
}

type Reply struct {
	Result int
}

type RegistryArgs struct {
	Service ServiceInfo
}