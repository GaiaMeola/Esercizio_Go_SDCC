package registry

import (
    "encoding/json"
    "net/http"
    "strings"
)

type Handler struct {
    store *Store
}

func NewHandler(store *Store) *Handler {
    return &Handler{store: store}
}

// Register aggiunge una nuova istanza
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
        return
    }

    var service Service
    if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
        http.Error(w, "JSON non valido", http.StatusBadRequest)
        return
    }

    if service.ID == "" || service.Name == "" {
        http.Error(w, "ID e Name obbligatori", http.StatusBadRequest)
        return
    }

    h.store.Register(service)
    w.WriteHeader(http.StatusCreated)
}

// Deregister rimuove una specifica istanza
// URL: /deregister/{serviceName}/{instanceID}
func (h *Handler) Deregister(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
        return
    }

    pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/deregister/"), "/")
    if len(pathParts) != 2 {
        http.Error(w, "URL deve essere /deregister/{serviceName}/{instanceID}", http.StatusBadRequest)
        return
    }

    serviceName := pathParts[0]
    instanceID := pathParts[1]

    if !h.store.Deregister(serviceName, instanceID) {
        http.Error(w, "Servizio o istanza non trovata", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// List restituisce tutte le istanze di tutti i servizi
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    services := h.store.List()
    json.NewEncoder(w).Encode(services)
}

// ListByName restituisce tutte le istanze di un servizio specifico
// URL: /services/{serviceName}
func (h *Handler) ListByName(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
        return
    }

    serviceName := strings.TrimPrefix(r.URL.Path, "/services/")
    if serviceName == "" {
        http.Error(w, "Nome servizio mancante", http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    services := h.store.ListByName(serviceName)
    json.NewEncoder(w).Encode(services)
}