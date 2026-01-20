#!/bin/bash

# Pulisce eventuali processi precedenti
killall registry server client 2>/dev/null

echo "--- Inizializzazione Ambiente ---"
mkdir -p state
echo "0" > state/counter.txt

echo "--- Avvio Service Registry ---"
go run registry/main.go &
sleep 2 # Attende che il registry sia pronto

echo "--- Avvio Server Replicati ---"
# Avvia Server 1 (Porta 8001, Peso 1)
go run server/main.go 8001 1 &
# Avvia Server 2 (Porta 8002, Peso 10)
go run server/main.go 8002 10 &
sleep 2

echo "--- Avvio Client ---"
echo "Premi CTRL+C per fermare tutto il sistema"
go run client/main.go

# Quando il client viene fermato, chiude anche i processi in background
trap "kill 0" EXIT