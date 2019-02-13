package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	tftypes "github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/pkg/api"
	rapi "github.com/threefoldtech/rivine/pkg/api"
)

type (
	// ERC20SyncingStatus contains a Ethereum syncing status.
	ERC20SyncingStatus struct {
		Status tftypes.ERC20SyncStatus `json:"status"`
	}
)

// RegisterERC20HTTPHandlers registers the (tfchain-specific) handlers for all ERC20 HTTP endpoints.
func RegisterERC20HTTPHandlers(router rapi.Router, erc20txValidator tftypes.ERC20TransactionValidator) {
	if erc20txValidator == nil {
		panic("no erc20Validator given")
	}
	if router == nil {
		panic("no router given")
	}

	// tfchain-specific endpoints

	router.GET("/erc20/downloader/status", NewERC20StatusHandler(erc20txValidator))
}

// NewERC20StatusHandler creates a handler to handle the API calls to /erc20/downloader/status.
func NewERC20StatusHandler(erc20txValidator tftypes.ERC20TransactionValidator) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		ERC20Status, err := erc20txValidator.GetStatus()
		if err != nil {
			api.WriteError(w, api.Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		api.WriteJSON(w, ERC20SyncingStatus{
			Status: *ERC20Status,
		})
	}
}
