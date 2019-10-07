package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/threefoldtech/rivine/types"
)

func (f *faucet) requestFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "" && r.URL.Path != "/" {
		http.Error(w, fmt.Errorf("%s is not a valid path", r.URL.Path).Error(), http.StatusNotFound)
		return
	}
	renderRequestTemplate(w, RequestBody{
		ChainName:    f.cts.ChainInfo.Name,
		ChainNetwork: f.cts.ChainInfo.NetworkName,
		CoinUnit:     f.cts.ChainInfo.CoinUnit,
	})
}

func (f *faucet) requestTokensHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] Parsing request token form")
	r.ParseForm()
	strUH := strings.Join(r.Form["uh"], "")
	var uh types.UnlockHash
	err := uh.LoadString(strUH)
	if err != nil {
		err = fmt.Errorf("invalid unlockhash %q: %v", strUH, err)
		renderRequestTemplate(w, RequestBody{
			ChainName:    f.cts.ChainInfo.Name,
			ChainNetwork: f.cts.ChainInfo.NetworkName,
			CoinUnit:     f.cts.ChainInfo.CoinUnit,
			Error:        err.Error(),
		})
		return
	}
	log.Println("[DEBUG] Requesting tokens for address", strUH)
	f.mu.Lock()
	defer f.mu.Unlock()
	txID, err := dripCoins(uh, f.coinsToGive)

	if err != nil {
		log.Println("[ERROR] Failed to drip coins:", err.Error())
		renderRequestTemplate(w, RequestBody{
			ChainName:    f.cts.ChainInfo.Name,
			ChainNetwork: f.cts.ChainInfo.NetworkName,
			CoinUnit:     f.cts.ChainInfo.CoinUnit,
			Error:        err.Error(),
		})
		return
	}
	renderCoinConfirmationTemplate(w, CoinConfirmationBody{
		ChainName:     f.cts.ChainInfo.Name,
		ChainNetwork:  f.cts.ChainInfo.NetworkName,
		CoinUnit:      f.cts.ChainInfo.CoinUnit,
		Address:       uh.String(),
		TransactionID: txID.String(),
	})
	log.Printf("[INFO] Sent %s tokens to %s\n", f.coinsToGive.String(), strUH)
}

func renderRequestTemplate(w http.ResponseWriter, body RequestBody) {
	err := requestTemplate.ExecuteTemplate(w, "request.html", body)
	if err != nil {
		log.Println("[ERROR] Failed to render template request.html:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderCoinConfirmationTemplate(w http.ResponseWriter, body CoinConfirmationBody) {
	err := coinConfirmationTemplate.ExecuteTemplate(w, "coinconfirmation.html", body)
	if err != nil {
		log.Println("[ERROR] Failed to render template coinconfirmation.html:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
