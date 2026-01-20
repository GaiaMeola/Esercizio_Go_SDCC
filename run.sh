#!/bin/bash

# Funzione per pulire i processi basandosi sulle porte nel file config.json
cleanup() {
    echo "--- Pulizia processi in corso... ---"
    
    # Estraiamo tutte le porte dal config.json per chiuderle forzatamente
    PORTS=$(grep -oP '"(port|registry_addr)": "\K[^"]+' config.json | grep -oP '\d+')
    
    for port in $PORTS; do
        fuser -k $port/tcp 2>/dev/null
    done
    
    # Chiudiamo tutti i processi Go e puliamo i file di log temporanei
    killall registry server client 2>/dev/null
    rm -f client_*.log
    sleep 1
}

echo "==============================================="
echo "   SDCC Project - Fully Configurable Runner    "
echo "==============================================="
echo "Cosa vuoi fare?"
echo "1) Avvia sistema (Registry + Servers + N Client da Config)"
echo "2) Solo pulizia (Reset)"
echo "3) Esci"
read -p "Scegli un'opzione [1-3]: " choice

case $choice in
    1)
        cleanup
        echo "--- Inizializzazione Stato ---"
        mkdir -p state
        echo "0" > state/counter.txt

        # 1. Estrazione e avvio Registry
        REG_ADDR=$(grep -oP '"registry_addr": "\K[^"]+' config.json)
        echo "--- Avvio Service Registry su $REG_ADDR ---"
        go run registry/main.go &
        sleep 2

        # 2. Estrazione e avvio Server Replicati
        PORTS=($(grep -oP '"port": "\K[^"]+' config.json))
        WEIGHTS=($(grep -oP '"weight": \K[^, }]+' config.json))

        echo "--- Avvio Server Replicati ---"
        for i in "${!PORTS[@]}"; do
            echo "Lancio Server sulla porta ${PORTS[$i]} con peso ${WEIGHTS[$i]}..."
            go run server/main.go ${PORTS[$i]} ${WEIGHTS[$i]} &
        done
        
        sleep 2

        # 3. Estrazione e avvio Client Multipli
        NUM_CLIENTS=$(grep "\"num_clients\":" config.json | grep -oP '\d+')
        if [ -z "$NUM_CLIENTS" ]; then NUM_CLIENTS=1; fi

        echo "--- Avvio di $NUM_CLIENTS Client in parallelo ---"
        # ... resto del ciclo for ...
        
        # Lanciamo N-1 client in background scrivendo i log su file
        for (( i=1; i<$NUM_CLIENTS; i++ )); do
            echo "Avvio Client $i in background (Log: client_$i.log)..."
            go run client/main.go > "client_$i.log" 2>&1 &
        done

        # L'ultimo client lo lanciamo in primo piano per vedere l'output a schermo
        echo "Avvio ultimo Client in primo piano. Premi CTRL+C per fermare tutto."
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