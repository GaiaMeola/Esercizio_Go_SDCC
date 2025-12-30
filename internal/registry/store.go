package registry

import "sync"

type Store struct {
    // Chiave: nome del servizio (o tipo)
    // Valore: mappa di istanze, chiave = ID istanza
    services map[string]map[string]Service
    mu       sync.RWMutex
}

func NewStore() *Store {
    return &Store{
        services: make(map[string]map[string]Service),
    }
}

// Register aggiunge una nuova istanza per un servizio
func (s *Store) Register(service Service) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, ok := s.services[service.Name]; !ok {
        s.services[service.Name] = make(map[string]Service)
    }

    s.services[service.Name][service.ID] = service
}

// Deregister rimuove una specifica istanza
func (s *Store) Deregister(serviceName, instanceID string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()

    instances, ok := s.services[serviceName]
    if !ok {
        return false
    }

    if _, exists := instances[instanceID]; !exists {
        return false
    }

    delete(instances, instanceID)

    // Se non ci sono più istanze, rimuoviamo la chiave del servizio
    if len(instances) == 0 {
        delete(s.services, serviceName)
    }

    return true
}

// List restituisce tutte le istanze di tutti i servizi
func (s *Store) List() []Service {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var all []Service
    for _, instances := range s.services {
        for _, svc := range instances {
            all = append(all, svc)
        }
    }
    return all
}

// ListByName restituisce tutte le istanze di un servizio specifico
func (s *Store) ListByName(serviceName string) []Service {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var instancesList []Service
    if instances, ok := s.services[serviceName]; ok {
        for _, svc := range instances {
            instancesList = append(instancesList, svc)
        }
    }
    return instancesList
}