package consensus

import (
	"errors"
	"net"
	"sync"
	"time"

	bolt "github.com/rivine/bbolt"
	"github.com/threefoldtech/rivine/build"
	"github.com/threefoldtech/rivine/crypto"
	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	"github.com/threefoldtech/rivine/types"
)

const (
	// minNumOutbound is the minimum number of outbound peers required before ibd
	// is confident we are synced.
	minNumOutbound = 3
)

var (
	// MaxCatchUpBlocks is the maxiumum number of blocks that can be given to
	// the consensus set in a single iteration during the initial blockchain
	// download.
	MaxCatchUpBlocks = func() types.BlockHeight {
		switch build.Release {
		case "dev":
			return 50
		case "testing":
			return 3
		default:
			if build.Release != "standard" {
				build.Severe("unrecognized build.Release")
			}
			return 10
		}
	}()
	// sendBlocksTimeout is the timeout for the SendBlocks RPC.
	sendBlocksTimeout = func() time.Duration {
		switch build.Release {
		case "dev":
			return 40 * time.Second
		case "testing":
			return 5 * time.Second
		default:
			if build.Release != "standard" {
				build.Severe("unrecognized build.Release")
			}
			return 5 * time.Minute
		}
	}()
	// ibdLoopDelay is the time that threadedInitialBlockchainDownload waits
	// between attempts to synchronize with the network if the last attempt
	// failed.
	ibdLoopDelay = func() time.Duration {
		switch build.Release {
		case "dev":
			return 1 * time.Second
		case "testing":
			return 100 * time.Millisecond
		default:
			if build.Release != "standard" {
				build.Severe("unrecognized build.Release")
			}
			return 10 * time.Second
		}
	}()

	errSendBlocksStalled = errors.New("SendBlocks RPC timed and never received any blocks")
)

// isTimeoutErr is a helper function that returns true if err was caused by a
// network timeout.
func isTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	// COMPATv1.3.0
	return (err.Error() == "Read timeout" || err.Error() == "Write timeout")
}

// blockHistory returns up to 32 block ids, starting with recent blocks and
// then proving exponentially increasingly less recent blocks. The genesis
// block is always included as the last block. This block history can be used
// to find a common parent that is reasonably recent, usually the most recent
// common parent is found, but always a common parent within a factor of 2 is
// found.
func blockHistory(tx *bolt.Tx) (blockIDs [32]types.BlockID) {
	height := blockHeight(tx)
	step := types.BlockHeight(1)
	// The final step is to include the genesis block, which is why the final
	// element is skipped during iteration.
	for i := 0; i < 31; i++ {
		// Include the next block.
		blockID, err := getPath(tx, height)
		if err != nil {
			build.Severe(err)
		}
		blockIDs[i] = blockID

		// Determine the height of the next block to include and then increase
		// the step size. The height must be decreased first to prevent
		// underflow.
		//
		// `i >= 9` means that the first 10 blocks will be included, and then
		// skipping will start.
		if i >= 9 {
			step *= 2
		}
		if height <= step {
			break
		}
		height -= step
	}
	// Include the genesis block as the last element
	blockID, err := getPath(tx, 0)
	if err != nil {
		build.Severe(err)
	}
	blockIDs[31] = blockID
	return blockIDs
}

// managedReceiveBlocks is the calling end of the SendBlocks RPC, without the
// threadgroup wrapping.
func (cs *ConsensusSet) managedReceiveBlocks(conn modules.PeerConn) (returnErr error) {
	// Set a deadline after which SendBlocks will timeout. During IBD, esepcially,
	// SendBlocks will timeout. This is by design so that IBD switches peers to
	// prevent any one peer from stalling IBD.
	err := conn.SetDeadline(time.Now().Add(sendBlocksTimeout))
	// Ignore errors returned by SetDeadline if the conn is a pipe in testing.
	// Pipes do not support Set{,Read,Write}Deadline and should only be used in
	// testing.
	if opErr, ok := err.(*net.OpError); ok && opErr.Op == "set" && opErr.Net == "pipe" && build.Release == "testing" {
		err = nil
	}
	if err != nil {
		return err
	}
	stalled := true
	defer func() {
		if stalled && isTimeoutErr(returnErr) {
			returnErr = errSendBlocksStalled
		}
	}()

	// Get blockIDs to send.
	var history [32]types.BlockID
	cs.mu.RLock()
	err = cs.db.View(func(tx *bolt.Tx) error {
		history = blockHistory(tx)
		return nil
	})
	cs.mu.RUnlock()
	if err != nil {
		return err
	}

	// Send the block ids.
	if err := siabin.WriteObject(conn, history); err != nil {
		return err
	}

	// Broadcast the last block accepted. This functionality is in a defer to
	// ensure that a block is always broadcast if any blocks are accepted. This
	// is to stop an attacker from preventing block broadcasts.
	chainExtended := false
	defer func() {
		cs.mu.RLock()
		synced := cs.synced
		cs.mu.RUnlock()
		if chainExtended && synced {
			// The last block received will be the current block since
			// managedAcceptBlock only returns nil if a block extends the longest chain.
			currentBlock := cs.managedCurrentBlock()
			peers := cs.gateway.Peers()
			go cs.gateway.Broadcast("RelayHeader", currentBlock.Header(), peers)
		}
	}()

	// Read blocks off of the wire and add them to the consensus set until
	// there are no more blocks available.
	moreAvailable := true
	for moreAvailable {
		// We need a check to see if we are stopping the loop, otherwise
		// we end up syncing the entire blockchain before exiting
		select {
		case <-cs.tg.StopChan():
			return nil
		default:
		}
		// Read a slice of blocks from the wire.
		var newBlocks []types.Block
		if err := siabin.ReadObject(conn, &newBlocks, uint64(MaxCatchUpBlocks)*cs.chainCts.BlockSizeLimit); err != nil {
			return err
		}
		if err := siabin.ReadObject(conn, &moreAvailable, 1); err != nil {
			return err
		}

		// Integrate the blocks into the consensus set.
		for _, block := range newBlocks {
			stalled = false
			// Call managedAcceptBlock instead of AcceptBlock so as not to broadcast
			// every block.
			acceptErr := cs.managedAcceptBlock(block)
			// Set a flag to indicate that we should broadcast the last block received.
			if acceptErr == nil {
				chainExtended = true
			}
			// ErrNonExtendingBlock must be ignored until headers-first block
			// sharing is implemented, block already in database should also be
			// ignored.
			if acceptErr == modules.ErrNonExtendingBlock || acceptErr == modules.ErrBlockKnown {
				acceptErr = nil
			}
			if acceptErr != nil {
				return acceptErr
			}
		}
	}
	return nil
}

// threadedReceiveBlocks is the calling end of the SendBlocks RPC.
func (cs *ConsensusSet) threadedReceiveBlocks(conn modules.PeerConn) error {
	err := cs.tg.Add()
	if err != nil {
		return err
	}
	defer cs.tg.Done()
	return cs.managedReceiveBlocks(conn)
}

// rpcSendBlocks is the receiving end of the SendBlocks RPC. It returns a
// sequential set of blocks based on the 32 input block IDs. The most recent
// known ID is used as the starting point, and up to 'MaxCatchUpBlocks' from
// that BlockHeight onwards are returned. It also sends a boolean indicating
// whether more blocks are available.
func (cs *ConsensusSet) rpcSendBlocks(conn modules.PeerConn) error {
	err := cs.tg.Add()
	if err != nil {
		return err
	}
	defer cs.tg.Done()

	// Read a list of blocks known to the requester and find the most recent
	// block from the current path.
	var knownBlocks [32]types.BlockID
	err = siabin.ReadObject(conn, &knownBlocks, 32*crypto.HashSize)
	if err != nil {
		return err
	}

	// Find the most recent block from knownBlocks in the current path.
	found := false
	var start types.BlockHeight
	var csHeight types.BlockHeight
	cs.mu.RLock()
	err = cs.db.View(func(tx *bolt.Tx) error {
		csHeight = blockHeight(tx)
		for _, id := range knownBlocks {
			pb, err := getBlockMap(tx, id)
			if err != nil {
				continue
			}
			pathID, err := getPath(tx, pb.Height)
			if err != nil {
				continue
			}
			if pathID != pb.Block.ID() {
				continue
			}
			if pb.Height == csHeight {
				break
			}
			found = true
			// Start from the child of the common block.
			start = pb.Height + 1
			break
		}
		return nil
	})
	cs.mu.RUnlock()
	if err != nil {
		return err
	}

	// If no matching blocks are found, or if the caller has all known blocks,
	// don't send any blocks.
	if !found {
		// Send 0 blocks.
		err = siabin.WriteObject(conn, []types.Block{})
		if err != nil {
			return err
		}
		// Indicate that no more blocks are available.
		return siabin.WriteObject(conn, false)
	}

	// Send the caller all of the blocks that they are missing.
	moreAvailable := true
	for moreAvailable {
		// Get the set of blocks to send.
		var blocks []types.Block
		cs.mu.RLock()
		err = cs.db.View(func(tx *bolt.Tx) error {
			height := blockHeight(tx)
			for i := start; i <= height && i < start+MaxCatchUpBlocks; i++ {
				id, err := getPath(tx, i)
				if err != nil {
					build.Severe(err)
				}
				pb, err := getBlockMap(tx, id)
				if err != nil {
					build.Severe(err)
				}
				blocks = append(blocks, pb.Block)
			}
			moreAvailable = start+MaxCatchUpBlocks <= height
			start += MaxCatchUpBlocks
			return nil
		})
		cs.mu.RUnlock()
		if err != nil {
			return err
		}

		// Send a set of blocks to the caller + a flag indicating whether more
		// are available.
		if err = siabin.WriteObject(conn, blocks); err != nil {
			return err
		}
		if err = siabin.WriteObject(conn, moreAvailable); err != nil {
			return err
		}
	}

	return nil
}

// threadedRPCRelayHeader is an RPC that accepts a block header from a peer.
func (cs *ConsensusSet) threadedRPCRelayHeader(conn modules.PeerConn) error {
	err := cs.tg.Add()
	if err != nil {
		return err
	}
	wg := new(sync.WaitGroup)
	defer func() {
		go func() {
			wg.Wait()
			cs.tg.Done()
		}()
	}()

	// Decode the block header from the connection.
	var h types.BlockHeader
	err = siabin.ReadObject(conn, &h, types.BlockHeaderSize)
	if err != nil {
		return err
	}

	// Start verification inside of a bolt View tx.
	cs.mu.RLock()
	err = cs.db.View(func(tx *bolt.Tx) error {
		// Do some relatively inexpensive checks to validate the header
		return cs.validateHeader(boltTxWrapper{tx}, h)
	})
	cs.mu.RUnlock()
	if err == errOrphan {
		// If the header is an orphan, try to find the parents. Call needs to
		// be made in a separate goroutine as execution requires calling an
		// exported gateway method - threadedRPCRelayHeader was likely called
		// from an exported gateway function.
		//
		// NOTE: In general this is bad design. Rather than recycling other
		// calls, the whole protocol should have been kept in a single RPC.
		// Because it is not, we have to do weird threading to prevent
		// deadlocks, and we also have to be concerned every time the code in
		// managedReceiveBlocks is adjusted.
		wg.Add(1)
		go func() {
			err := cs.gateway.RPC(conn.RPCAddr(), "SendBlocks", cs.managedReceiveBlocks)
			if err != nil {
				cs.log.Debugln("WARN: failed to get parents of orphan header:", err)
			}
			wg.Done()
		}()
		return nil
	} else if err != nil {
		return err
	}

	// If the header is valid and extends the heaviest chain, fetch the
	// corresponding block. Call needs to be made in a separate goroutine
	// because an exported call to the gateway is used, which is a deadlock
	// risk given that rpcRelayHeader is called from the gateway.
	//
	// NOTE: In general this is bad design. Rather than recycling other calls,
	// the whole protocol should have been kept in a single RPC. Because it is
	// not, we have to do weird threading to prevent deadlocks, and we also
	// have to be concerned every time the code in managedReceiveBlock is
	// adjusted.
	wg.Add(1)
	go func() {
		err = cs.gateway.RPC(conn.RPCAddr(), "SendBlk", cs.managedReceiveBlock(h.ID()))
		if err != nil {
			cs.log.Debugln("WARN: failed to get header's corresponding block:", err)
		}
		wg.Done()
	}()
	return nil
}

// rpcSendBlk is an RPC that sends the requested block to the requesting peer.
func (cs *ConsensusSet) rpcSendBlk(conn modules.PeerConn) error {
	err := cs.tg.Add()
	if err != nil {
		return err
	}
	defer cs.tg.Done()

	// Decode the block id from the connection.
	var id types.BlockID
	err = siabin.ReadObject(conn, &id, crypto.HashSize)
	if err != nil {
		return err
	}
	// Lookup the corresponding block.
	var b types.Block
	cs.mu.RLock()
	err = cs.db.View(func(tx *bolt.Tx) error {
		pb, err := getBlockMap(tx, id)
		if err != nil {
			return err
		}
		b = pb.Block
		return nil
	})
	cs.mu.RUnlock()
	if err != nil {
		return err
	}
	// Encode and send the block to the caller.
	err = siabin.WriteObject(conn, b)
	if err != nil {
		return err
	}
	return nil
}

// managedReceiveBlock takes a block id and returns an RPCFunc that requests that
// block and then calls AcceptBlock on it. The returned function should be used
// as the calling end of the SendBlk RPC. Note that although the function
// itself does not do any locking, it is still prefixed with "threaded" because
// the function it returns calls the exported method AcceptBlock.
func (cs *ConsensusSet) managedReceiveBlock(id types.BlockID) modules.RPCFunc {
	return func(conn modules.PeerConn) error {
		if err := siabin.WriteObject(conn, id); err != nil {
			return err
		}
		var block types.Block
		if err := siabin.ReadObject(conn, &block, cs.chainCts.BlockSizeLimit); err != nil {
			return err
		}
		if err := cs.managedAcceptBlock(block); err != nil {
			return err
		}
		cs.managedBroadcastBlock(block)
		return nil
	}
}

// threadedInitialBlockchainDownload performs the IBD on outbound peers. Blocks
// are downloaded from one peer at a time in 5 minute intervals, so as to
// prevent any one peer from significantly slowing down IBD.
//
// NOTE: IBD will succeed right now when each peer has a different blockchain.
// The height and the block id of the remote peers' current blocks are not
// checked to be the same.
func (cs *ConsensusSet) threadedInitialBlockchainDownload() error {
	// The consensus set will not recognize IBD as complete until it has enough
	// peers. After the deadline though, it will recognize the blockchain
	// download as complete even with only one outbound peer synced. This deadline is helpful
	// to local-net setups, where a machine will frequently only have one peer
	// (and that peer will be another machine on the same local network, but
	// within the local network at least one peer is connected to the broad
	// network).
	maxIBDWaitTime := time.Duration(cs.chainCts.BlockFrequency) * time.Second
	numOutboundSynced := 0
	numOutboundNotSynced := 0

	// keep track of our initial block height
	getHeight := func() (height types.BlockHeight) {
		_ = cs.db.View(func(tx *bolt.Tx) error {
			height = blockHeight(tx)
			return nil
		})
		return
	}
	height := getHeight()
	lastReceiveTime := time.Now()

	for {
		numOutboundSynced = 0
		numOutboundNotSynced = 0
		for _, p := range cs.gateway.Peers() {
			// We only sync on outbound peers at first to make IBD less susceptible to
			// fast-mining and other attacks, as outbound peers are more difficult to
			// manipulate.
			if p.Inbound {
				continue
			}

			// Put the rest of the iteration inside of a thread group.
			err := func() error {
				err := cs.tg.Add()
				if err != nil {
					return err
				}
				defer cs.tg.Done()

				// Request blocks from the peer. The error returned will only be
				// 'nil' if there are no more blocks to receive.
				err = cs.gateway.RPC(p.NetAddress, "SendBlocks", cs.managedReceiveBlocks)
				if err == nil {
					numOutboundSynced++

					// check if we moved up our block height
					currentHeight := getHeight()
					if currentHeight > height {
						height = currentHeight
						lastReceiveTime = time.Now()
					}

					// In this case, 'return nil' is equivalent to skipping to
					// the next iteration of the loop.
					return nil
				}

				numOutboundNotSynced++
				if isTimeoutErr(err) {
					cs.log.Printf("WARN: disconnecting from peer %v because IBD failed: %v", p.NetAddress, err)
					// Disconnect if there is an unexpected error (not a timeout). This
					// includes errSendBlocksStalled.
					//
					// We disconnect so that these peers are removed from gateway.Peers() and
					// do not prevent us from marking ourselves as fully synced.
					err := cs.gateway.Disconnect(p.NetAddress)
					if err != nil {
						cs.log.Printf("WARN: disconnecting from peer %v failed: %v", p.NetAddress, err)
					}
				}

				cs.log.Printf("WARN: managedReceiveBlocks (via SendBlocks RPC) has failed with an error: %v", err)
				return nil
			}()
			if err != nil {
				return err
			}
		}

		// The consensus set is not considered synced until a majority of
		// outbound peers say that we are synced. If less than 10 minutes have
		// passed, a minimum of 'minNumOutbound' peers must say that we are
		// synced, otherwise a 1 vs 0 majority is sufficient.
		//
		// This scheme is used to prevent malicious peers from being able to
		// barricade the sync'd status of the consensus set, and to make sure
		// that consensus sets behind a firewall with only one peer
		// (potentially a local peer) are still able to eventually conclude
		// that they have syncrhonized. Miners and hosts will often have setups
		// beind a firewall where there is a single node with many peers and
		// then the rest of the nodes only have a few peers.
		if numOutboundSynced > numOutboundNotSynced && numOutboundSynced >= minNumOutbound {
			cs.log.Printf("INFO: Stopping IBD, sufficient amount of outbound peers (%d) are in sync", numOutboundSynced)
			break
		}

		// if we didn't receive anything in a sufficient amount of time,
		// and at least one peer is synced, we will continue as well.
		if numOutboundSynced >= 1 && time.Since(lastReceiveTime) >= maxIBDWaitTime {
			cs.log.Printf(
				"INFO: Stopping IBD, only %d synced outbound peers, but no more blocks received in %v",
				numOutboundSynced, time.Since(lastReceiveTime))
			break
		}

		// Sleep so we don't hammer the network with SendBlock requests.
		time.Sleep(ibdLoopDelay)
	}

	cs.log.Printf("INFO: IBD done, synced with %v peers", numOutboundSynced)
	return nil
}

// Synced returns true if the consensus set is synced with the network.
func (cs *ConsensusSet) Synced() bool {
	err := cs.tg.Add()
	if err != nil {
		return false
	}
	defer cs.tg.Done()
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.synced
}
