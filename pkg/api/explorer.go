package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/threefoldfoundation/tfchain/pkg/persist"
	"github.com/threefoldfoundation/tfchain/pkg/types"
	"github.com/threefoldtech/rivine/build"
	"github.com/threefoldtech/rivine/crypto"
	"github.com/threefoldtech/rivine/modules"
	rapi "github.com/threefoldtech/rivine/pkg/api"
	rtypes "github.com/threefoldtech/rivine/types"
)

// ExplorerHashGET wraps around the default rivine ExplorerHashGET type,
// as to add the optional ERC20 address to it, for UnlockHash requests,
// which have an ERC20 address attached to it.
type ExplorerHashGET struct {
	rapi.ExplorerHashGET
	ERC20Info *ExplorerHashERC20Info `json:"erc20info,omitempty"`
}

// ExplorerHashERC20Info contains all ERC20 related info as part of an UnlockHash-typed ExplorerHashGET request.
type ExplorerHashERC20Info struct {
	TFTAddress    rtypes.UnlockHash  `json:"tftaddress,omitempty"`
	ERC20Address  types.ERC20Address `json:"erc20address,omitempty"`
	Confirmations uint64             `json:"confirmations"`
}

// RegisterExplorerHTTPHandlers registers the (tfchain-specific) handlers for all Explorer HTTP endpoints.
func RegisterExplorerHTTPHandlers(router rapi.Router, cs modules.ConsensusSet, explorer modules.Explorer, tpool modules.TransactionPool, txdb *persist.TransactionDB) {
	if cs == nil {
		panic("no ConsensusSet API given")
	}
	if explorer == nil {
		panic("no Explorer API given")
	}
	if txdb == nil {
		panic("no TransactiondB API given")
	}
	if router == nil {
		panic("no router given")
	}

	// rivine endpoints

	router.GET("/explorer", rapi.NewExplorerRootHandler(explorer))
	router.GET("/explorer/blocks/:height", rapi.NewExplorerBlocksHandler(cs, explorer))
	router.GET("/explorer/stats/history", rapi.NewExplorerHistoryStatsHandler(explorer))
	router.GET("/explorer/stats/range", rapi.NewExplorerRangeStatsHandler(explorer))
	router.GET("/explorer/constants", rapi.NewExplorerConstantsHandler(explorer))

	// tfchain-specific endpoints

	router.GET("/explorer/mintcondition", NewTransactionDBGetActiveMintConditionHandler(txdb))
	router.GET("/explorer/mintcondition/:height", NewTransactionDBGetMintConditionAtHandler(txdb))

	router.GET("/explorer/3bot/:id", NewTransactionDBGetRecordForIDHandler(txdb))
	router.GET("/explorer/whois/3bot/:name", NewTransactionDBGetRecordForNameHandler(txdb))
	router.GET("/explorer/3bot/:id/transactions", NewTransactionDBGetBotTransactionsHandler(txdb))

	router.GET("/explorer/erc20/addresses/:address", NewTransactionDBGetERC20RelatedAddressHandler(txdb))
	router.GET("/explorer/erc20/transactions/:txid", NewTransactionDBGetERC20TransactionID(txdb))

	// tfchain rivine-overwritten endpoints

	router.GET("/explorer/hashes/:hash", NewExplorerHashHandler(explorer, cs, tpool, txdb))
}

// NewExplorerHashHandler creates a handler to handle GET requests to /explorer/hash/:hash.
func NewExplorerHashHandler(explorer modules.Explorer, cs modules.ConsensusSet, tpool modules.TransactionPool, txdb *persist.TransactionDB) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// Scan the hash as a hash. If that fails, try scanning the hash as an
		// address.
		hash, err := rapi.ScanHash(ps.ByName("hash"))
		if err != nil {
			var erc20Address types.ERC20Address
			var erc20AddressUnconfirmed bool

			hstr := ps.ByName("hash")
			addr, err := rapi.ScanAddress(hstr)
			var found bool
			if err != nil {
				if len(hstr) != types.ERC20AddressLength*2 {
					rapi.WriteError(w, rapi.Error{Message: err.Error()}, http.StatusBadRequest)
					return
				}

				// decode the ERC20 address
				err = erc20Address.LoadString(hstr)
				if err != nil {
					rapi.WriteError(w, rapi.Error{Message: "invalid address: only valid TFT/ERC20 addresses are accepted"}, http.StatusBadRequest)
					return
				}

				// get the TFT address using the ERC20 address
				addr, found, err = txdb.GetTFTAddressForERC20Address(erc20Address)
				if err != nil {
					rapi.WriteError(w, rapi.Error{Message: "invalid ERC20 address: " + err.Error()}, http.StatusBadRequest)
					return
				}
				if !found && tpool != nil {
					addr, found, err = getERC20AddressRegInfoFromTxPool(tpool, erc20Address)
					if err == nil && !found {
						err = errors.New("address not found")
					}
					if err != nil {
						rapi.WriteError(w, rapi.Error{Message: "invalid ERC20 address: " + err.Error()}, http.StatusBadRequest)
						return
					}
					erc20AddressUnconfirmed = true
				}
			} else {
				// try to get the ERC20 Address
				// ignore error: is not critical
				erc20Address, found, err = txdb.GetERC20AddressForTFTAddress(addr)
				if err != nil {
					if build.DEBUG {
						log.Printf("error while fetching ERC20 address for TFT Address %v: %v", addr, err)
					}
					erc20AddressUnconfirmed = true
				} else if !found {
					// try to find it in the TxPool if possible
					err = nil
					if tpool != nil {
						erc20Address = types.ERC20AddressFromUnlockHash(addr)
						_, found, err = getERC20AddressRegInfoFromTxPool(tpool, erc20Address)
					}
					if err == nil && !found {
						err = errors.New("address not found")
					}
					if err != nil && build.DEBUG {
						log.Printf("error while fetching ERC20 address for TFT Address %v: %v", addr, err)
					}
					erc20AddressUnconfirmed = true
				}
			}

			// Try the hash as an unlock hash. Unlock hash is checked last because
			// unlock hashes do not have collision-free guarantees. Someone can create
			// an unlock hash that collides with another object id. They will not be
			// able to use the unlock hash, but they can disrupt the explorer. This is
			// handled by checking the unlock hash last. Anyone intentionally creating
			// a colliding unlock hash (such a collision can only happen if done
			// intentionally) will be unable to find their unlock hash in the
			// blockchain through the explorer hash lookup.
			var (
				txns   []rapi.ExplorerTransaction
				blocks []rapi.ExplorerBlock
			)
			if txids := explorer.UnlockHash(addr); len(txids) != 0 {
				// parse the optional filters for the unlockhash request
				var filters rapi.TransactionSetFilters
				if str := req.FormValue("minheight"); str != "" {
					n, err := strconv.ParseUint(str, 10, 64)
					if err != nil {
						rapi.WriteError(w, rapi.Error{Message: "invalid minheight filter: " + err.Error()}, http.StatusBadRequest)
						return
					}
					filters.MinBlockHeight = rtypes.BlockHeight(n)
				}
				// build the transaction set for all transactions for the given unlock hash,
				// taking into account the given filters
				txns, blocks = rapi.BuildTransactionSet(explorer, txids, filters)
			}
			txns = append(txns, getUnconfirmedTransactions(explorer, tpool, addr)...)
			multiSigAddresses := explorer.MultiSigAddresses(addr)
			if len(txns) != 0 || len(blocks) != 0 || len(multiSigAddresses) != 0 || erc20Address != (types.ERC20Address{}) {
				resp := ExplorerHashGET{
					ExplorerHashGET: rapi.ExplorerHashGET{
						HashType:          rapi.HashTypeUnlockHashStr,
						Blocks:            blocks,
						Transactions:      txns,
						MultiSigAddresses: multiSigAddresses,
					},
				}
				if erc20Address != (types.ERC20Address{}) {
					resp.ERC20Info = &ExplorerHashERC20Info{
						TFTAddress:   addr,
						ERC20Address: erc20Address,
					}
					// If not unconfirmed, ensure to add the confirmation height
					fmt.Println(erc20AddressUnconfirmed, txns)
					if !erc20AddressUnconfirmed {
						curHeight := cs.Height()
						regHeight := curHeight
						for _, txn := range txns {
							if txn.RawTransaction.Version != types.TransactionVersionERC20AddressRegistration {
								continue
							}
							regtxn, _ := types.ERC20AddressRegistrationTransactionFromTransaction(txn.RawTransaction)
							if types.ERC20AddressFromUnlockHash(rtypes.NewPubKeyUnlockHash(regtxn.PublicKey)) == erc20Address {
								regHeight = txn.Height
								break
							}
						}
						if curHeight >= regHeight {
							resp.ERC20Info.Confirmations = uint64((curHeight - regHeight) + 1)
						}
					}
				}
				rapi.WriteJSON(w, resp)
				return
			}

			// Hash not found, return an error.
			rapi.WriteError(w, rapi.Error{Message: "no transactions or blocks found for given address"}, http.StatusNoContent)
			return
		}

		// TODO: lookups on the zero hash are too expensive to allow. Need a
		// better way to handle this case.
		if hash == (crypto.Hash{}) {
			rapi.WriteError(w, rapi.Error{Message: "can't lookup the empty unlock hash"}, http.StatusBadRequest)
			return
		}

		// Try the hash as a block id.
		block, height, exists := explorer.Block(rtypes.BlockID(hash))
		if exists {
			rapi.WriteJSON(w, ExplorerHashGET{
				ExplorerHashGET: rapi.ExplorerHashGET{
					HashType: rapi.HashTypeBlockIDStr,
					Block:    rapi.BuildExplorerBlock(explorer, height, block),
				},
			})
			return
		}

		// Try the hash as a transaction id.
		block, height, exists = explorer.Transaction(rtypes.TransactionID(hash))
		if exists {
			var txn rtypes.Transaction
			for _, t := range block.Transactions {
				if t.ID() == rtypes.TransactionID(hash) {
					txn = t
				}
			}
			rapi.WriteJSON(w, ExplorerHashGET{
				ExplorerHashGET: rapi.ExplorerHashGET{
					HashType:    rapi.HashTypeTransactionIDStr,
					Transaction: rapi.BuildExplorerTransaction(explorer, height, block.ID(), txn),
				},
			})
			return
		}

		// Try the hash as a siacoin output id.
		txids := explorer.CoinOutputID(rtypes.CoinOutputID(hash))
		if len(txids) != 0 {
			txns, blocks := rapi.BuildTransactionSet(explorer, txids, rapi.TransactionSetFilters{})
			rapi.WriteJSON(w, ExplorerHashGET{
				ExplorerHashGET: rapi.ExplorerHashGET{
					HashType:     rapi.HashTypeCoinOutputIDStr,
					Blocks:       blocks,
					Transactions: txns,
				},
			})
			return
		}

		// Try the hash as a siafund output id.
		txids = explorer.BlockStakeOutputID(rtypes.BlockStakeOutputID(hash))
		if len(txids) != 0 {
			txns, blocks := rapi.BuildTransactionSet(explorer, txids, rapi.TransactionSetFilters{})
			rapi.WriteJSON(w, ExplorerHashGET{
				ExplorerHashGET: rapi.ExplorerHashGET{
					HashType:     rapi.HashTypeBlockStakeOutputIDStr,
					Blocks:       blocks,
					Transactions: txns,
				},
			})
			return
		}

		// if the transaction pool is available, try to use it
		if tpool != nil {
			// Try the hash as a transactionID in the transaction pool
			txn, err := tpool.Transaction(rtypes.TransactionID(hash))
			if err == nil {
				rapi.WriteJSON(w, ExplorerHashGET{
					ExplorerHashGET: rapi.ExplorerHashGET{
						HashType:    rapi.HashTypeTransactionIDStr,
						Transaction: rapi.BuildExplorerTransaction(explorer, 0, rtypes.BlockID{}, txn),
						Unconfirmed: true,
					},
				})
				return
			}
			if err != modules.ErrTransactionNotFound {
				rapi.WriteError(w, rapi.Error{
					Message: "error during call to /explorer/hash: failed to get txn from transaction pool: " + err.Error()},
					http.StatusInternalServerError)
				return
			}
		}

		// Hash not found, return an error.
		rapi.WriteError(w, rapi.Error{Message: "unrecognized hash used as input to /explorer/hash"}, http.StatusBadRequest)
	}
}

func getERC20AddressRegInfoFromTxPool(tpool modules.TransactionPool, erc20Addr types.ERC20Address) (rtypes.UnlockHash, bool, error) {
	for _, txn := range tpool.TransactionList() {
		if txn.Version != types.TransactionVersionERC20AddressRegistration {
			continue
		}
		regtxn, err := types.ERC20AddressRegistrationTransactionFromTransaction(txn)
		if err != nil {
			return rtypes.UnlockHash{}, false, err
		}
		uh := rtypes.NewPubKeyUnlockHash(regtxn.PublicKey)
		if types.ERC20AddressFromUnlockHash(uh) == erc20Addr {
			return uh, true, nil
		}
	}
	return rtypes.UnlockHash{}, false, nil
}

// getUnconfirmedTransactions returns a list of all transactions which are unconfirmed and related to the given unlock hash from the transactionpool
func getUnconfirmedTransactions(explorer modules.Explorer, tpool modules.TransactionPool, addr rtypes.UnlockHash) []rapi.ExplorerTransaction {
	if tpool == nil {
		return nil
	}
	relatedTxns := []rtypes.Transaction{}
	unconfirmedTxns := tpool.TransactionList()
	// make a list of potential unspend coin outputs
	potentiallySpentCoinOutputs := map[rtypes.CoinOutputID]rtypes.CoinOutput{}
	for _, txn := range unconfirmedTxns {
		for idx, sco := range txn.CoinOutputs {
			potentiallySpentCoinOutputs[txn.CoinOutputID(uint64(idx))] = sco
		}
	}
	// go through all unconfirmed transactions
	for _, txn := range unconfirmedTxns {
		related := false
		// Check if any coin output is related to the hash we currently have
		for _, co := range txn.CoinOutputs {
			if co.Condition.UnlockHash() == addr {
				related = true
				relatedTxns = append(relatedTxns, txn)
				break
			}
		}
		if related {
			continue
		}
		// Check if any blockstake output is related
		for _, bso := range txn.BlockStakeOutputs {
			if bso.Condition.UnlockHash() == addr {
				related = true
				relatedTxns = append(relatedTxns, txn)
				break
			}
		}
		if related {
			continue
		}
		// Check the coin inputs
		for _, ci := range txn.CoinInputs {
			// check if related to an unconfirmed coin output
			if sco, ok := potentiallySpentCoinOutputs[ci.ParentID]; ok && sco.Condition.UnlockHash() == addr {
				// mark related, add tx and stop this coin input loop
				related = true
				relatedTxns = append(relatedTxns, txn)
				break
			}
			// check if related to a confirmed coin output
			co, _ := explorer.CoinOutput(ci.ParentID)
			if co.Condition.UnlockHash() == addr {
				related = true
				relatedTxns = append(relatedTxns, txn)
				break
			}
		}
		if related {
			continue
		}
		// Check blockstake inputs
		for _, bsi := range txn.BlockStakeInputs {
			bsi, _ := explorer.BlockStakeOutput(bsi.ParentID)
			if bsi.Condition.UnlockHash() == addr {
				related = true
				relatedTxns = append(relatedTxns, txn)
				break
			}
		}
	}
	explorerTxns := make([]rapi.ExplorerTransaction, len(relatedTxns))
	for i := range relatedTxns {
		relatedTxn := relatedTxns[i]
		spentCoinOutputs := map[rtypes.CoinOutputID]rtypes.CoinOutput{}
		for _, sci := range relatedTxn.CoinInputs {
			// add unconfirmed coin output
			if sco, ok := potentiallySpentCoinOutputs[sci.ParentID]; ok {
				spentCoinOutputs[sci.ParentID] = sco
				continue
			}
			// add confirmed coin output
			sco, exists := explorer.CoinOutput(sci.ParentID)
			if build.DEBUG && !exists {
				panic("could not find corresponding coin output")
			}
			spentCoinOutputs[sci.ParentID] = sco
		}
		explorerTxns[i] = buildExplorerTransactionWithMappedCoinOutputs(explorer, 0, rtypes.BlockID{}, relatedTxn, spentCoinOutputs)
		explorerTxns[i].Unconfirmed = true
	}
	return explorerTxns
}

func buildExplorerTransactionWithMappedCoinOutputs(explorer modules.Explorer, height rtypes.BlockHeight, parent rtypes.BlockID, txn rtypes.Transaction, spentCoinOutputs map[rtypes.CoinOutputID]rtypes.CoinOutput) (et rapi.ExplorerTransaction) {
	// Get the header information for the transaction.
	et.ID = txn.ID()
	et.Height = height
	et.Parent = parent
	et.RawTransaction = txn

	// Add the siacoin outputs that correspond with each siacoin input.
	for _, sci := range txn.CoinInputs {
		sco, ok := spentCoinOutputs[sci.ParentID]
		if build.DEBUG && !ok {
			panic("could not find corresponding coin output")
		}
		et.CoinInputOutputs = append(et.CoinInputOutputs, rapi.ExplorerCoinOutput{
			CoinOutput: sco,
			UnlockHash: sco.Condition.UnlockHash(),
		})
	}

	for i, co := range txn.CoinOutputs {
		et.CoinOutputIDs = append(et.CoinOutputIDs, txn.CoinOutputID(uint64(i)))
		et.CoinOutputUnlockHashes = append(et.CoinOutputUnlockHashes, co.Condition.UnlockHash())
	}

	// Add the siafund outputs that correspond to each siacoin input.
	for _, sci := range txn.BlockStakeInputs {
		sco, exists := explorer.BlockStakeOutput(sci.ParentID)
		if build.DEBUG && !exists {
			panic("could not find corresponding blockstake output")
		}
		et.BlockStakeInputOutputs = append(et.BlockStakeInputOutputs, rapi.ExplorerBlockStakeOutput{
			BlockStakeOutput: sco,
			UnlockHash:       sco.Condition.UnlockHash(),
		})
	}

	for i, bso := range txn.BlockStakeOutputs {
		et.BlockStakeOutputIDs = append(et.BlockStakeOutputIDs, txn.BlockStakeOutputID(uint64(i)))
		et.BlockStakeOutputUnlockHashes = append(et.BlockStakeOutputUnlockHashes, bso.Condition.UnlockHash())
	}

	return et
}
