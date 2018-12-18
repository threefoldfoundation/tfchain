package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/threefoldfoundation/tfchain/pkg/persist"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"

	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/types"

	"github.com/julienschmidt/httprouter"
)

type (
	// TransactionDBGetMintCondition contains a requested mint condition,
	// either the current active one active for the given blockheight or lower.
	TransactionDBGetMintCondition struct {
		MintCondition types.UnlockConditionProxy `json:"mintcondition"`
	}

	// TransactionDBGetBotRecord contains a requested bot record.
	TransactionDBGetBotRecord struct {
		Record tftypes.BotRecord `json:"record"`
	}

	// TransactionDBGetBotTransactions contains the requested identifiers
	// of transactions for a specific bot.
	TransactionDBGetBotTransactions struct {
		Identifiers []types.TransactionID `json:"ids"`
	}

	// TransactionDBGetERC20RelatedAddress contains the requested ERC20-related addresses.
	TransactionDBGetERC20RelatedAddress struct {
		TFTAddress   types.UnlockHash     `json:"tftaddress"`
		ERC20Address tftypes.ERC20Address `json:"erc20address"`
	}

	// TransactionDBGetERC20TransactionID contains the requested info found for the given ERC20 TransactionID.
	TransactionDBGetERC20TransactionID struct {
		ERC20TransaxtionID   tftypes.ERC20TransactionID `json:"er20txid"`
		TfchainTransactionID types.TransactionID        `json:"tfttxid"`
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

	router.GET("/consensus/3bot/:id", NewTransactionDBGetRecordForIDHandler(txdb))
	router.GET("/consensus/whois/3bot/:name", NewTransactionDBGetRecordForNameHandler(txdb))
	router.GET("/consensus/3bot/:id/transactions", NewTransactionDBGetBotTransactionsHandler(txdb))

	router.GET("/consensus/erc20/addresses/:address", NewTransactionDBGetERC20RelatedAddressHandler(txdb))
	router.GET("/consensus/erc20/transactions/:txid", NewTransactionDBGetERC20TransactionID(txdb))
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

// NewTransactionDBGetRecordForIDHandler creates a handler to handle the API calls to /transactiondb/3bot/:id.
func NewTransactionDBGetRecordForIDHandler(txdb *persist.TransactionDB) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var (
			err    error
			record *tftypes.BotRecord
		)
		idStr := ps.ByName("id")
		var id tftypes.BotID
		err = id.LoadString(idStr)
		if err == nil {
			// interpret it as a BotID
			record, err = txdb.GetRecordForID(tftypes.BotID(id))
		} else {
			// interpret it as a PublicKey
			var pubKey types.PublicKey
			err = pubKey.LoadString(idStr)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Errorf("id has to be a valid PublicKey or BotID: %v", err).Error()},
					http.StatusBadRequest)
				return
			}
			record, err = txdb.GetRecordForKey(pubKey)
		}
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		api.WriteJSON(w, TransactionDBGetBotRecord{
			Record: *record,
		})
	}
}

// NewTransactionDBGetRecordForNameHandler creates a handler to handle the API calls to /transactiondb/whois/3bot/:name.
func NewTransactionDBGetRecordForNameHandler(txdb *persist.TransactionDB) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var name tftypes.BotName
		err := name.LoadString(ps.ByName("name"))
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Errorf("invalid botname: %v", err).Error()},
				http.StatusInternalServerError)
			return
		}
		record, err := txdb.GetRecordForName(name)
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		api.WriteJSON(w, TransactionDBGetBotRecord{
			Record: *record,
		})
	}
}

// NewTransactionDBGetBotTransactionsHandler creates a handler to handle the API calls to /transactiondb/3bot/:id/transactions.
func NewTransactionDBGetBotTransactionsHandler(txdb *persist.TransactionDB) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		idStr := ps.ByName("id")
		var id tftypes.BotID
		err := id.LoadString(idStr)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Errorf("id has to be a valid BotID: %v", err).Error()},
				http.StatusBadRequest)
			return
		}
		ids, err := txdb.GetBotTransactionIdentifiers(id)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Errorf("failed to get transactions for BotID: %v", err).Error()},
				http.StatusInternalServerError)
			return
		}
		api.WriteJSON(w, TransactionDBGetBotTransactions{
			Identifiers: ids,
		})
	}
}

// NewTransactionDBGetERC20RelatedAddressHandler creates a handler to handle the API calls to /transactiondb/erc20/addresses/:address.
func NewTransactionDBGetERC20RelatedAddressHandler(txdb *persist.TransactionDB) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		addressStr := ps.ByName("address")

		var (
			err  error
			resp TransactionDBGetERC20RelatedAddress
		)
		if len(addressStr) == tftypes.ERC20AddressLength*2 {
			err = resp.ERC20Address.LoadString(addressStr)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Sprintf("invalid ERC20 address given: %v", err)}, http.StatusBadRequest)
				return
			}
			resp.TFTAddress, err = txdb.GetTFTAddressForERC20Address(resp.ERC20Address)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Sprintf("error while fetching TFT Address: %v", err)}, http.StatusInternalServerError)
				return
			}
		} else {
			err = resp.TFTAddress.LoadString(addressStr)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Sprintf("invalid (TFT) address given: %v", err)}, http.StatusBadRequest)
				return
			}
			resp.ERC20Address, err = txdb.GetERC20AddressForTFTAddress(resp.TFTAddress)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Sprintf("error while fetching ERC20 Address: %v", err)}, http.StatusInternalServerError)
				return
			}
		}
		api.WriteJSON(w, resp)
	}
}

// NewTransactionDBGetERC20TransactionID creates a handler to handle the API calls to /transactiondb/erc20/transactions/:txid.
func NewTransactionDBGetERC20TransactionID(txdb *persist.TransactionDB) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		txidStr := ps.ByName("txid")
		var txid tftypes.ERC20TransactionID
		err := txid.LoadString(txidStr)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Sprintf("invalid ERC20 TransactionID given: %v", err)}, http.StatusBadRequest)
			return
		}

		tfttxid, err := txdb.GetTFTTransactionIDForERC20TransactionID(txid)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Sprintf("error while fetching info linked to ERC20 TransactionID: %v", err)}, http.StatusInternalServerError)
			return
		}

		api.WriteJSON(w, TransactionDBGetERC20TransactionID{
			ERC20TransaxtionID:   txid,
			TfchainTransactionID: tfttxid,
		})
	}
}
