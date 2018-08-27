package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/threefoldfoundation/tfchain/pkg/persist"

	"github.com/rivine/rivine/pkg/api"
	"github.com/rivine/rivine/types"

	"github.com/julienschmidt/httprouter"
)

type (
	// TransactionDBGetMintCondition contains a requested mint condition,
	// either the current active one active for the given blockheight or lower.
	TransactionDBGetMintCondition struct {
		MintCondition types.UnlockConditionProxy `json:"mintcondition"`
	}
)

// RegisterTransactionDBHTTPHandlers registers the handlers for all TransactionDB HTTP endpoints.
func RegisterTransactionDBHTTPHandlers(router api.Router, txdb *persist.TransactionDB) {
	if txdb == nil {
		panic("no transaction DB given")
	}
	if router == nil {
		panic("no httprouter Router given")
	}

	router.GET("/consensus/mintcondition", NewTransactionDBGetActiveMintConditionHandler(txdb))
	router.GET("/consensus/mintcondition/:height", NewTransactionDBGetMintConditionAtHandler(txdb))
}

// NewTransactionDBGetActiveMintConditionHandler creates a handler to handle the API calls to /transactiondb/mintcondition.
func NewTransactionDBGetActiveMintConditionHandler(txdb *persist.TransactionDB) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		mintCondition, err := txdb.GetActiveMintCondition()
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		api.WriteJSON(w, TransactionDBGetMintCondition{
			MintCondition: mintCondition,
		})
	}
}

// NewTransactionDBGetMintConditionAtHandler creates a handler to handle the API calls to /transactiondb/mintcondition/:height.
func NewTransactionDBGetMintConditionAtHandler(txdb *persist.TransactionDB) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		heightStr := ps.ByName("height")
		height, err := strconv.ParseUint(heightStr, 10, 64)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Sprintf("invalid block height given: %v", err)}, http.StatusBadRequest)
			return
		}
		mintCondition, err := txdb.GetMintConditionAt(types.BlockHeight(height))
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		api.WriteJSON(w, TransactionDBGetMintCondition{
			MintCondition: mintCondition,
		})
	}
}
