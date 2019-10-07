package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/threefoldtech/rivine/types"
)

func (f *faucet) requestCoins(w http.ResponseWriter, r *http.Request) {
	addCorsHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body := struct {
		Address types.UnlockHash `json:"address"`
		Amount  uint64           `json:"amount"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] Requesting coins (%s) through API\n", body.Address.String())

	f.mu.Lock()
	defer f.mu.Unlock()

	var txID types.TransactionID
	if body.Amount == 0 {
		txID, err = dripCoins(body.Address, f.coinsToGive)
	} else {
		// If there is an amount requested, use the provided amount
		txID, err = dripCoins(body.Address, f.cts.OneCoin.Mul64(body.Amount))
	}

	if err != nil {
		log.Println("[ERROR] Failed to drip coins:", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		TxID types.TransactionID `json:"txid"`
	}{TxID: txID})
}

func addCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
