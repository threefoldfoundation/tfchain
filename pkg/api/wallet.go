package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/types"
)

type (
	// WalletFundCoins is the resulting object that is returned,
	// to be used by a client to fund a transaction of any type.
	WalletFundCoins struct {
		CoinInputs       []types.CoinInput `json:"coininputs"`
		RefundCoinOutput *types.CoinOutput `json:"refund"`
	}

	// WalletPublicKeyGET contains a public key returned by a GET call to
	// /wallet/publickey.
	WalletPublicKeyGET struct {
		PublicKey types.PublicKey `json:"publickey"`
	}
)

// RegisterWalletHTTPHandlers registers the (tfchain-specific) handlers for all Wallet HTTP endpoints.
func RegisterWalletHTTPHandlers(router api.Router, wallet modules.Wallet, requiredPassword string) {
	if wallet == nil {
		panic("no wallet API given")
	}
	if router == nil {
		panic("no httprouter Router given")
	}

	router.GET("/wallet/publickey", api.RequirePasswordHandler(NewWalletGetPublicKeyHandler(wallet), requiredPassword))
	router.GET("/wallet/fund/coins", api.RequirePasswordHandler(NewWalletFundCoinsHandler(wallet), requiredPassword))
}

// NewWalletFundCoinsHandler creates a handler to handle the API calls to /wallet/fund/coins?amount=.
// While it might be handy for other use cases, it is needed for 3bot registration
func NewWalletFundCoinsHandler(wallet modules.Wallet) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		q := req.URL.Query()
		// parse the amount
		amountStr := q.Get("amount")
		if amountStr == "" || amountStr == "0" {
			api.WriteError(w, api.Error{Message: "an amount has to be specified, greater than 0"}, http.StatusBadRequest)
			return
		}
		var amount types.Currency
		err := amount.LoadString(amountStr)
		if err != nil {
			api.WriteError(w, api.Error{Message: "invalid amount given: " + err.Error()}, http.StatusBadRequest)
			return
		}

		// start a transaction and fund the requested amount
		txbuilder := wallet.StartTransaction()
		err = txbuilder.FundCoins(amount)
		if err != nil {
			api.WriteError(w, api.Error{Message: "failed to fund the requested coins: " + err.Error()}, http.StatusInternalServerError)
			return
		}

		// build the dummy Txn, as to view the Txn
		txn, _ := txbuilder.View()
		// defer drop the Txn
		defer txbuilder.Drop()

		// compose the result object and validate it
		result := WalletFundCoins{CoinInputs: txn.CoinInputs}
		if len(result.CoinInputs) == 0 {
			api.WriteError(w, api.Error{Message: "no coin inputs could be generated"}, http.StatusInternalServerError)
			return
		}
		switch len(txn.CoinOutputs) {
		case 0:
			// ignore, valid, but nothing to do
		case 1:
			// add as refund
			result.RefundCoinOutput = &txn.CoinOutputs[0]
		case 2:
			api.WriteError(w, api.Error{Message: "more than 2 coin outputs were generated, while maximum 1 was expected"}, http.StatusInternalServerError)
			return
		}
		// all good, return the resulting object
		api.WriteJSON(w, result)
	}
}

// NewWalletGetPublicKeyHandler creates a handler to handle API calls to /wallet/publickey.
// While it might be handy for other use cases, it is needed for 3bot.
func NewWalletGetPublicKeyHandler(wallet modules.Wallet) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		unlockHash, err := wallet.NextAddress()
		if err != nil {
			api.WriteError(w, api.Error{Message: "error after call to /wallet/publickey: " + err.Error()}, walletErrorToHTTPStatus(err))
			return
		}
		pk, _, err := wallet.GetKey(unlockHash)
		if err != nil {
			api.WriteError(w, api.Error{Message: "failed to fetch newly created public key: " + err.Error()}, http.StatusInternalServerError)
			return
		}
		api.WriteJSON(w, WalletPublicKeyGET{PublicKey: pk})
	}
}

func walletErrorToHTTPStatus(err error) int {
	if err == modules.ErrLockedWallet {
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}
