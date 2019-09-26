package http

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/threefoldtech/rivine/pkg/api"
	"github.com/threefoldtech/rivine/types"

	erc20types "github.com/threefoldtech/rivine-extension-erc20/types"
)

type (
	// ERC20SyncingStatus contains a Ethereum syncing status.
	ERC20SyncingStatus struct {
		Status erc20types.ERC20SyncStatus `json:"status"`
	}
	// ERC20BalanceInformation contains the Bridge Contract's balance info
	ERC20BalanceInformation struct {
		BalanceInfo erc20types.ERC20BalanceInfo `json:"balanceinfo"`
	}
	// GetERC20RelatedAddress contains the requested ERC20-related addresses.
	GetERC20RelatedAddress struct {
		TFTAddress   types.UnlockHash        `json:"tftaddress"`
		ERC20Address erc20types.ERC20Address `json:"erc20address"`
	}

	// GetERC20TransactionID contains the requested info found for the given ERC20 TransactionID.
	GetERC20TransactionID struct {
		ERC20TransaxtionID   erc20types.ERC20Hash `json:"er20txid"`
		TfchainTransactionID types.TransactionID  `json:"tfttxid"`
	}
)

// RegisterERC20HTTPHandlers registers the (tfchain-specific) handlers for all ERC20 HTTP endpoints.
func RegisterERC20HTTPHandlers(router api.Router, erc20InfoAPI erc20types.ERC20InfoAPI) {
	if erc20InfoAPI == nil {
		panic("no erc20InfoApi given")
	}
	if router == nil {
		panic("no router given")
	}

	// tfchain-specific endpoints

	router.GET("/erc20/downloader/status", NewERC20StatusHandler(erc20InfoAPI))
	router.GET("/erc20/account/balance", newERC20BalanceHandler(erc20InfoAPI))

}

// NewERC20StatusHandler creates a handler to handle the API calls to /erc20/downloader/status.
func NewERC20StatusHandler(erc20InfoAPI erc20types.ERC20InfoAPI) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		ERC20Status, err := erc20InfoAPI.GetStatus()
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		api.WriteJSON(w, ERC20SyncingStatus{
			Status: *ERC20Status,
		})
	}
}

// newERC20BalanceHandler creates a handler to handle the API calls to /erc20/account/balance.
func newERC20BalanceHandler(erc20InfoAPI erc20types.ERC20InfoAPI) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		ERC20BalanceInfo, err := erc20InfoAPI.GetBalanceInfo()
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		api.WriteJSON(w, ERC20BalanceInformation{
			BalanceInfo: *ERC20BalanceInfo,
		})
	}
}

// RegisterConsensusHTTPHandlers registers the 3Bot handlers for all consensus HTTP endpoints.
func RegisterConsensusHTTPHandlers(router api.Router, erc20Registry erc20types.ERC20Registry) {
	if erc20Registry == nil {
		panic("no ERC20Registry API given")
	}
	if router == nil {
		panic("no httprouter Router given")
	}

	router.GET("/consensus/erc20/addresses/:address", NewGetERC20RelatedAddressHandler(erc20Registry))
	router.GET("/consensus/erc20/transactions/:txid", NewGetERC20TransactionID(erc20Registry))
}

// RegisterExplorerHTTPHandlers registers the 3Bot handlers for all explorer HTTP endpoints.
func RegisterExplorerHTTPHandlers(router api.Router, erc20Registry erc20types.ERC20Registry) {
	if erc20Registry == nil {
		panic("no ERC20Registry API given")
	}
	if router == nil {
		panic("no httprouter Router given")
	}

	router.GET("/explorer/erc20/addresses/:address", NewGetERC20RelatedAddressHandler(erc20Registry))
	router.GET("/explorer/erc20/transactions/:txid", NewGetERC20TransactionID(erc20Registry))
}

// NewGetERC20RelatedAddressHandler creates a handler to handle the API calls to /transactiondb/erc20/addresses/:address.
func NewGetERC20RelatedAddressHandler(erc20Registry erc20types.ERC20Registry) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		addressStr := ps.ByName("address")

		var (
			err   error
			found bool
			resp  GetERC20RelatedAddress
		)
		if len(addressStr) == erc20types.ERC20AddressLength*2 {
			err = resp.ERC20Address.LoadString(addressStr)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Sprintf("invalid ERC20 address given: %v", err)}, http.StatusBadRequest)
				return
			}
			resp.TFTAddress, found, err = erc20Registry.GetTFTAddressForERC20Address(resp.ERC20Address)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Sprintf("error while fetching TFT Address: %v", err)}, http.StatusInternalServerError)
				return
			}
			if !found {
				api.WriteError(w, api.Error{Message: "error while fetching TFT Address: address not found"}, http.StatusNoContent)
				return
			}
		} else {
			err = resp.TFTAddress.LoadString(addressStr)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Sprintf("invalid (TFT) address given: %v", err)}, http.StatusBadRequest)
				return
			}
			resp.ERC20Address, found, err = erc20Registry.GetERC20AddressForTFTAddress(resp.TFTAddress)
			if err != nil {
				api.WriteError(w, api.Error{Message: fmt.Sprintf("error while fetching ERC20 Address: %v", err)}, http.StatusInternalServerError)
				return
			}
			if !found {
				api.WriteError(w, api.Error{Message: "error while fetching ERC20 Address: address not found"}, http.StatusNoContent)
				return
			}
		}
		api.WriteJSON(w, resp)
	}
}

// NewGetERC20TransactionID creates a handler to handle the API calls to /transactiondb/erc20/transactions/:txid.
func NewGetERC20TransactionID(erc20Registry erc20types.ERC20Registry) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		txidStr := ps.ByName("txid")
		var txid erc20types.ERC20Hash
		err := txid.LoadString(txidStr)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Sprintf("invalid ERC20 TransactionID given: %v", err)}, http.StatusBadRequest)
			return
		}

		tfttxid, found, err := erc20Registry.GetTFTTransactionIDForERC20TransactionID(txid)
		if err != nil {
			api.WriteError(w, api.Error{Message: fmt.Sprintf("error while fetching info linked to ERC20 TransactionID: %v", err)}, http.StatusInternalServerError)
			return
		}
		if !found {
			api.WriteError(w, api.Error{Message: "error while fetching info linked to ERC20 TransactionID: ID not found"}, http.StatusNoContent)
			return
		}

		api.WriteJSON(w, GetERC20TransactionID{
			ERC20TransaxtionID:   txid,
			TfchainTransactionID: tfttxid,
		})
	}
}
