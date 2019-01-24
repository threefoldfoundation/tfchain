package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/spf13/cobra"
	"github.com/threefoldfoundation/tfchain/pkg/eth/erc20"
)

const (
	precision = 1000000000.
)

// Commands defines the CLI Commands for the tfchain client.
type Commands struct {
	EthNetworkName string

	// eth port for light client
	EthPort uint16

	// eth bootnodes
	EthBootNodes []string

	// eth account flags
	accJSON string
	accPass string

	EthLog          int
	ContractAddress string

	RootPersistentDir string

	// port to make the server listen on
	ServerAddr uint16
}

type demoExchange struct {
	contract *erc20.BridgeContract
}

// GetAddress returns the loaded address of the demo exchange
func (de *demoExchange) GetAddress(w http.ResponseWriter, r *http.Request) {
	addr, err := de.contract.AccountAddress()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := struct {
			Error error `json:"error"`
		}{
			Error: err,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := struct {
		Address common.Address `json:"address"`
	}{
		Address: addr,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (de *demoExchange) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := de.contract.EthBalance()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := struct {
			Error error `json:"error"`
		}{
			Error: err,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := struct {
		Balance *big.Int `json:"balance"`
	}{
		Balance: balance,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (de *demoExchange) GetTokenBalance(w http.ResponseWriter, r *http.Request) {
	addr, err := de.contract.AccountAddress()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := struct {
			Error error `json:"error"`
		}{
			Error: err,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	balance, err := de.contract.TokenBalance(addr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := struct {
			Error error `json:"error"`
		}{
			Error: err,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := struct {
		Balance *big.Int `json:"balance"`
	}{
		Balance: balance,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (de *demoExchange) Withdraw(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Address common.Address `json:"address"`
		Amount  *big.Int       `json:"amount"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Info("Failed to decode withdraw body", "err", err)
		w.WriteHeader(http.StatusBadRequest)

		res := struct {
			Error error `json:"error"`
		}{
			Error: err,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(body)

	err = de.contract.TransferFunds(body.Address, body.Amount)
	if err != nil {
		log.Info("Failed to transfer tokens", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		res := struct {
			Error error `json:"error"`
		}{
			Error: err,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Root represents the root command,
func (cmd *Commands) Root(_ *cobra.Command, args []string) error {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(cmd.EthLog), log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	closeChan := make(chan struct{})

	contract, err := erc20.NewBridgeContract(strings.ToLower(cmd.EthNetworkName), cmd.EthBootNodes, cmd.ContractAddress, int(cmd.EthPort), cmd.accJSON, cmd.accPass, cmd.RootPersistentDir, closeChan)
	if err != nil {
		log.Error("Failed to create contract bindings", "err", err)
		return err
	}

	de := &demoExchange{contract: contract}

	server := http.Server{Addr: ":" + strconv.Itoa(int(cmd.ServerAddr))}
	mux := http.NewServeMux()
	mux.HandleFunc("/balance", de.GetBalance)
	mux.HandleFunc("/address", de.GetAddress)
	mux.HandleFunc("/tokenbalance", de.GetTokenBalance)
	mux.HandleFunc("/withdraw", de.Withdraw)

	mux.Handle("/", http.FileServer(http.Dir("./webpages")))

	server.Handler = mux

	if err := server.ListenAndServe(); err != nil {
		log.Error("Server error", "err", err)
		return err
	}

	return nil
}

func main() {
	cmd := new(Commands)

	// define commands
	cmdRoot := &cobra.Command{
		Use:          "demo_exchange",
		Short:        "start the demo_exchange",
		Long:         `start the demo_exchange`,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE:         cmd.Root,
	}

	// define flags
	cmdRoot.Flags().StringVarP(
		&cmd.RootPersistentDir,
		"persistent-directory", "d",
		"demo_exchange",
		"location of the root directory used to store persistent data",
	)

	// bridge flags
	cmdRoot.Flags().StringVar(
		&cmd.EthNetworkName,
		"ethnetwork", "ropsten",
		"The ethereum network, {rinkeby, ropsten}",
	)
	cmdRoot.Flags().Uint16Var(
		&cmd.EthPort,
		"ethport", 30303,
		"port for the ethereum deamon",
	)
	cmdRoot.Flags().StringSliceVar(
		&cmd.EthBootNodes,
		"ethbootnodes", nil,
		"Override the default ethereum bootnodes, a comma seperated list of enode URLs (enode://pubkey1@ip1:port1)",
	)

	// bridge account
	cmdRoot.Flags().StringVar(
		&cmd.accJSON,
		"account-json", "",
		"the path to an account file. If set, the specified account will be loaded",
	)

	cmdRoot.Flags().StringVar(
		&cmd.accPass,
		"account-password", "",
		"Password for the bridge account",
	)

	cmdRoot.Flags().IntVarP(
		&cmd.EthLog,
		"ethereum-log-lvl", "e", 3,
		"Log lvl for the ethereum logger",
	)

	cmdRoot.Flags().StringVar(
		&cmd.ContractAddress,
		"contract-address", "",
		"Use a custom contract",
	)

	cmdRoot.Flags().Uint16Var(
		&cmd.ServerAddr,
		"port", 8080,
		"Port to serve the demo web page on",
	)

	// execute logic
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
