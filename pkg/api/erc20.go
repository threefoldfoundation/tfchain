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
	// ERC20BalanceInformation contains the Bridge Contract's balance info
	ERC20BalanceInformation struct {
		BalanceInfo tftypes.ERC20BalanceInfo `json:"balanceinformation"`
	}
)

// RegisterERC20HTTPHandlers registers the (tfchain-specific) handlers for all ERC20 HTTP endpoints.
func RegisterERC20HTTPHandlers(router rapi.Router, erc20InfoAPI tftypes.ERC20InfoAPI) {
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
func NewERC20StatusHandler(erc20InfoAPI tftypes.ERC20InfoAPI) httprouter.Handle {
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
func newERC20BalanceHandler(erc20InfoAPI tftypes.ERC20InfoAPI) httprouter.Handle {
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
