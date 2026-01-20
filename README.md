# Progetto Sistemi Distribuiti: RPC Service Registry & Load Balancing

**Studente:** Gaia Meola  
**Corso:** Sistemi Distribuiti e Cloud Computing (SDCC)  
**Anno Accademico:** 2025/2026

---

## 1. Obiettivo del Progetto
Sviluppo di un'applicazione distribuita in **Go** basata su **RPC** che implementa il pattern **Client-Side Service Discovery**. Il sistema permette la gestione dinamica di microservizi replicati, garantendo il bilanciamento del carico e la consistenza dello stato tra le repliche.



## 2. Architettura del Sistema
Il sistema si basa su tre componenti principali che interagiscono secondo il modello *Service-Oriented*:

* **Service Registry (Service Discovery Server)**: Nodo centrale responsabile del monitoraggio dei server attivi. Gestisce le operazioni di `Register`, `Deregister` e fornisce la lista dei provider ai client tramite `Lookup`.
* **Service Providers (Replicated Servers)**: Istanze multiple che espongono:
    * **Servizio Stateless**: Una funzione di somma (`Add`) che non dipende da stati precedenti.
    * **Servizio Stateful**: Un contatore globale condiviso. La consistenza tra le repliche è garantita tramite uno **Shared File Storage** (`state/counter.txt`), simulando un database esterno come suggerito dai requisiti delle slide (Pag. 2).
* **Service Consumer (Client)**: Implementa la logica di scoperta e consumo dei servizi. Utilizza una **Cache locale** per ottimizzare le performance e ridurre il traffico verso il Registry.

## 3. Funzionalità Tecniche Implementate

### Bilanciamento del Carico (Load Balancing)
Il client implementa algoritmi di selezione dinamica per distribuire le richieste RPC:
* **Weighted Load Balancing**: Distribuisce il carico in base al "peso" (capacità computazionale) dichiarato da ogni server nel file di configurazione.
* **Fault Tolerance**: In caso di fallimento di un server, il client invalida la cache locale e forza una nuova `Lookup` per garantire la continuità del servizio.



### Meccanismo di Caching
* **TTL (Time-To-Live)**: La lista dei server è mantenuta in memoria per un intervallo configurabile (es. 15s) per minimizzare l'overhead di rete.
* **Aggiornamento Dinamico**: La cache viene aggiornata alla scadenza del TTL o immediatamente in caso di errore di connessione.

### Configurazione Centralizzata
Parametri come porte, indirizzi, numero di client e pesi dei server sono isolati nel file `config.json`, rendendo il sistema scalabile senza modifiche al codice sorgente.

---

## 4. Struttura del Progetto
```text
service-registry-go/
├── common/           # Strutture dati, definizioni RPC e logica JSON
├── registry/         # Implementazione del Service Registry
├── server/           # Logica del Service Provider (Stateless + Stateful)
├── client/           # Logica del Service Consumer, Load Balancer e Cache
├── state/            # Directory per la persistenza dello stato condiviso
├── config.json       # File di configurazione centralizzato
└── run.sh            # Script di orchestrazione e gestione processi

5. Guida all'Esecuzione
Prerequisiti

    Go 1.20 o superiore.

    Ambiente Linux/WSL2 per l'esecuzione dello script di automazione.

Avvio Interattivo

Il progetto include uno script di orchestrazione che automatizza la pulizia delle porte e l'avvio sincronizzato dei componenti:

    Posizionarsi nella root del progetto: cd service-registry-go

    Rendere lo script eseguibile: chmod +x run.sh

    Eseguire lo script: ./run.sh

    Selezionare l'opzione 1 per avviare l'intero ecosistema.

Verifica dei Requisiti

    Consistenza: Il contatore globale incrementa in modo coerente indipendentemente dal server selezionato dal Load Balancer.

    Discovery: Spegnendo un server (CTRL+C), il Registry lo rimuove e il Client smette di inviargli richieste dopo l'aggiornamento della cache.