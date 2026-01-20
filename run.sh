#!/bin/bash

# Funzione per pulire i processi basandosi sulle porte nel file config.json
cleanup() {
    echo "--- Pulizia processi in corso... ---"
    
    # Estraiamo tutte le porte dal config.json per chiuderle forzatamente
    # Legge la porta del registry e tutte le porte dei server
    PORTS=$(grep -oP '"(port|registry_addr)": "\K[^"]+' config.json | grep -oP '\d+')
    
    for port in $PORTS; do
        fuser -k $port/tcp 2>/dev/null
    done
    
    killall registry server client 2>/dev/null
    sleep 1
}

echo "==============================================="
echo "   SDCC Project - Configurable Runner          "
echo "==============================================="
echo "Cosa vuoi fare?"
echo "1) Avvia sistema (Registry + Servers da Config + Client)"
echo "2) Solo pulizia (Reset)"
echo "3) Esci"
read -p "Scegli un'opzione [1-3]: " choice

case $choice in
    1)
        cleanup
        echo "--- Inizializzazione Stato ---"
        mkdir -p state
        echo "0" > state/counter.txt

        # Estraiamo l'indirizzo del registry (es. localhost:5000)
        REG_ADDR=$(grep -oP '"registry_addr": "\K[^"]+' config.json)
        echo "--- Avvio Service Registry su $REG_ADDR ---"
        go run registry/main.go &
        sleep 2

        echo "--- Avvio Server Replicati ---"
        # Estraiamo liste di porte e pesi usando grep e trasformandoli in array
        PORTS=($(grep -oP '"port": "\K[^"]+' config.json))
        WEIGHTS=($(grep -oP '"weight": \K[^, }]+' config.json))

        for i in "${!PORTS[@]}"; do
            echo "Lancio Server sulla porta ${PORTS[$i]} con peso ${WEIGHTS[$i]}..."
            go run server/main.go ${PORTS[$i]} ${WEIGHTS[$i]} &
        done
        
        sleep 2
        echo "--- Avvio Client ---"
        echo "Configurazione caricata. Premi CTRL+C per fermare tutto."
        go run client/main.go
        ;;
    2)
        cleanup
        echo "Sistema pulito."
        ;;
    3)
        exit 0
        ;;
    *)
        echo "Opzione non valida."
        ;;
esac

# Chiude tutto automaticamente quando chiudi lo script
trap "cleanup; exit" INT TERM EXIT