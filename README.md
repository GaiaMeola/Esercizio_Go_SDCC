# Progetto Sistemi Distribuiti: RPC Service Registry & Load Balancing

**Studente:** Gaia Meola
**Corso:** Sistemi Distribuiti e Cloud Computing (SDCC)
**Anno Accademico:** 2025/2026

## 1. Obiettivo del Progetto
Realizzazione di un'applicazione distribuita in **Go** che utilizza **RPC** per implementare un sistema di Service Discovery lato client. Il sistema permette a più server di registrarsi presso un Registry centrale e a più client di consumare servizi (stateless e stateful) con bilanciamento del carico.

## 2. Architettura del Sistema
L'architettura segue il pattern **Client-Side Service Discovery**:

* **Service Registry**: Un nodo centrale che mantiene la lista dei server attivi (IP:Porta e Peso).
* **Replicated Servers**: Più istanze che espongono:
    * **Servizio Stateless**: Una funzione di calcolo (es. Somma) che non dipende da chiamate precedenti.
    * **Servizio Stateful**: Un contatore globale la cui consistenza è garantita tra tutte le repliche tramite un sistema di archiviazione comune (File System locale).
* **Client**: Interroga il Registry, salva i server in una **Cache locale** e applica algoritmi di **Load Balancing**.



---

## 3. Funzionalità Implementate

### Service Discovery & Lifecycle
- **Self-Registration**: All'avvio, ogni server comunica al Registry il proprio indirizzo e il proprio "peso".
- **Deregistration**: Allo spegnimento, il server informa il Registry per essere rimosso dalla lista.
- **Lookup**: Il client ottiene la lista aggiornata dei server disponibili.

### Load Balancing & Caching (Lato Client)
- **Caching Dinamica**: Il client non interroga il Registry a ogni chiamata, ma mantiene una cache locale con un TTL (Time-To-Live).
- **Invalidazione Automatica**: Se un server non risponde, il client invalida la cache e richiede una nuova lista al Registry.
- **Algoritmi di Bilanciamento**:
    - **Round Robin**: Selezione ciclica dei server.
    - **Weighted**: Selezione basata sul peso del server (maggiore è il peso, più richieste riceve).

### Consistenza dello Stato (Stateful Service)
Per garantire che il servizio stateful sia consistente tra le repliche senza l'uso di Docker Volumes, il sistema utilizza un file condiviso (`state/counter.txt`) protetto da meccanismi di sincronizzazione, simulando un database esterno.

---

## 4. Struttura del Progetto
```text
service-registry-go/
├── common/           # Strutture dati condivise e definizioni RPC
├── registry/         # Codice del Service Registry
├── server/           # Logica dei server (Stateless + Stateful)
├── client/           # Logica del client, Load Balancer e Cache
└── state/            # Directory locale per la persistenza dello stato