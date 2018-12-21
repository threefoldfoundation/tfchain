package erc20

import (
	"os"
	"path/filepath"

	"github.com/threefoldtech/rivine/modules"
	"github.com/threefoldtech/rivine/persist"
	"github.com/threefoldtech/rivine/types"
)

const (
	settingsFile = "bridge.json"
)

var (
	settingsMetadata = persist.Metadata{
		Header:  "Bridge Settings",
		Version: "0.0.1",
	}
)

type (
	// persist contains all of the persistent miner data.
	persistence struct {
		RecentChange modules.ConsensusChangeID
		Height       types.BlockHeight
	}
)

// initSettings loads the settings file if it exists and creates it if it
// doesn't.
func (bridge *Bridge) initSettings() error {
	filename := filepath.Join(bridge.persistDir, settingsFile)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return bridge.save()
	} else if err != nil {
		return err
	}
	return bridge.load()
}

// initPersist initializes the persistence of the bridge.
func (bridge *Bridge) initPersist() error {
	// Create the miner directory.
	err := os.MkdirAll(bridge.persistDir, 0700)
	if err != nil {
		return err
	}

	return bridge.initSettings()
}

// load loads the bridge persistence from disk.
func (bridge *Bridge) load() error {
	return persist.LoadJSON(settingsMetadata, &bridge.persist, filepath.Join(bridge.persistDir, settingsFile))
}

// save saves the bridge persistence to disk.
func (bridge *Bridge) save() error {
	return persist.SaveJSON(settingsMetadata, bridge.persist, filepath.Join(bridge.persistDir, settingsFile))
}
