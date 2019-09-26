package api

import (
	"fmt"
	"net/http"

	tbtypes "github.com/threefoldfoundation/tfchain/extensions/threebot/types"

	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/types"

	"github.com/julienschmidt/httprouter"
)

type (
	// GetBotRecord contains a requested bot record.
	GetBotRecord struct {
		Record tbtypes.BotRecord `json:"record"`
	}

	// GetBotTransactions contains the requested identifiers
	// of transactions for a specific bot.
	GetBotTransactions struct {
		Identifiers []types.TransactionID `json:"ids"`
	}
)

// RegisterConsensusHTTPHandlers registers the 3Bot handlers for all consensus HTTP endpoints.
func RegisterConsensusHTTPHandlers(router api.Router, tbRegistry tbtypes.BotRecordReadRegistry) {
	if tbRegistry == nil {
		panic("no BotRecordReadRegistry API given")
	}
	if router == nil {
		panic("no httprouter Router given")
	}

	router.GET("/consensus/3bot/:id", NewGetRecordForIDHandler(tbRegistry))
	router.GET("/consensus/whois/3bot/:name", NewGetRecordForNameHandler(tbRegistry))
	router.GET("/consensus/3bot/:id/transactions", NewGetBotTransactionsHandler(tbRegistry))
}

// RegisterExplorerHTTPHandlers registers the 3Bot handlers for all explorer HTTP endpoints.
func RegisterExplorerHTTPHandlers(router api.Router, tbRegistry tbtypes.BotRecordReadRegistry) {
	if tbRegistry == nil {
		panic("no BotRecordReadRegistry API given")
	}
	if router == nil {
		panic("no httprouter Router given")
	}

	router.GET("/explorer/3bot/:id", NewGetRecordForIDHandler(tbRegistry))
	router.GET("/explorer/whois/3bot/:name", NewGetRecordForNameHandler(tbRegistry))
	router.GET("/explorer/3bot/:id/transactions", NewGetBotTransactionsHandler(tbRegistry))
}

// NewGetRecordForIDHandler creates a handler to handle the API calls to /transactiondb/3bot/:id.
func NewGetRecordForIDHandler(tbRegistry tbtypes.BotRecordReadRegistry) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var (
			err    error
			record *tbtypes.BotRecord
		)
		idStr := ps.ByName("id")
		var id tbtypes.BotID
		err = id.LoadString(idStr)
		if err == nil {
			// interpret it as a BotID
			record, err = tbRegistry.GetRecordForID(tbtypes.BotID(id))
		} else {
			// interpret it as a PublicKey
			var pubKey types.PublicKey
			err = pubKey.LoadString(idStr)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Errorf("id has to be a valid PublicKey or BotID: %v", err).Error()},
					http.StatusBadRequest)
				return
			}
			record, err = tbRegistry.GetRecordForKey(pubKey)
		}
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, threeBotErrorAsHTTPStatusCode(err))
			return
		}
		api.WriteJSON(w, GetBotRecord{
			Record: *record,
		})
	}
}

// NewGetRecordForNameHandler creates a handler to handle the API calls to /transactiondb/whois/3bot/:name.
func NewGetRecordForNameHandler(tbRegistry tbtypes.BotRecordReadRegistry) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var name tbtypes.BotName
		err := name.LoadString(ps.ByName("name"))
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Errorf("invalid botname: %v", err).Error()},
				http.StatusInternalServerError)
			return
		}
		record, err := tbRegistry.GetRecordForName(name)
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, threeBotErrorAsHTTPStatusCode(err))
			return
		}
		api.WriteJSON(w, GetBotRecord{
			Record: *record,
		})
	}
}

// NewGetBotTransactionsHandler creates a handler to handle the API calls to /transactiondb/3bot/:id/transactions.
func NewGetBotTransactionsHandler(tbRegistry tbtypes.BotRecordReadRegistry) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		idStr := ps.ByName("id")
		var id tbtypes.BotID
		err := id.LoadString(idStr)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Errorf("id has to be a valid BotID: %v", err).Error()},
				http.StatusBadRequest)
			return
		}
		ids, err := tbRegistry.GetBotTransactionIdentifiers(id)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Errorf("failed to get transactions for BotID: %v", err).Error()},
				threeBotErrorAsHTTPStatusCode(err))
			return
		}
		api.WriteJSON(w, GetBotTransactions{
			Identifiers: ids,
		})
	}
}

// threeBotErrorAsHTTPStatusCode converts a 3bot error to an http status code.
// if it is not an applicable 3bot error, an internal server error code is returned
func threeBotErrorAsHTTPStatusCode(err error) int {
	switch err {
	case tbtypes.ErrBotNotFound, tbtypes.ErrBotNameNotFound, tbtypes.ErrBotKeyNotFound:
		return http.StatusNotFound
	case tbtypes.ErrBotNameExpired:
		return http.StatusPaymentRequired
	default:
		return http.StatusInternalServerError
	}
}
