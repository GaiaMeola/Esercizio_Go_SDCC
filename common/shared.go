package common

// ServiceInfo contiene i dettagli che il Registry condivide con i Client
type ServiceInfo struct {
	Addr   string // Indirizzo (es. "localhost:8001")
	Weight int    // Peso per il Load Balancing Weighted
}

// ArgsStateless: Per un servizio semplice come la somma (Stateless)
type ArgsStateless struct {
	A, B int
}

// ArgsStateful: Per il contatore condiviso (Stateful)
type ArgsStateful struct {
	Value int // Quanto vogliamo incrementare
}

// Reply: Risposta generica che restituisce un intero (risultato o valore contatore)
type Reply struct {
	Result int
}

// RegistryArgs: Usato dai Server per registrarsi/deregistrarsi
type RegistryArgs struct {
	Service ServiceInfo
}