package wallet

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/threefoldtech/rivine/modules"
)

const (
	// tfhcainDir is the root storage location for the tfchain files
	tfchaindDir = ".tfchain"
	// walletsSubDir is the location where the wallet files are stored
	walletsSubDir = "light-wallets"
	// walletFileName is the name of the wallet save file
	walletFileName = "wallet.json"
)

type (
	walletPersist struct {
		Seed       modules.Seed `json:"seed"`
		KeysToLoad uint64       `json:"keys_to_load"`
	}
)

func walletExists(name string) (bool, error) {
	_, err := os.Stat(Dir(name))
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func save(wallet *Wallet) error {
	data := walletPersist{
		Seed:       wallet.seed,
		KeysToLoad: uint64(len(wallet.keys)),
	}
	err := os.MkdirAll(Dir(wallet.name), 0777)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(Dir(wallet.name), walletFileName))
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(data)
}

func load(name string) (walletPersist, error) {
	exists, err := walletExists(name)
	if err != nil {
		return walletPersist{}, err
	}
	if !exists {
		return walletPersist{}, ErrNoSuchWallet
	}
	w := walletPersist{}
	file, err := os.Open(filepath.Join(Dir(name), walletFileName))
	if err != nil {
		return w, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&w)
	return w, err
}

// UserHomeDir gets the home directory of the current user
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// PersistDir is the directory in which the light wallet data is stored
func PersistDir() string {
	return filepath.Join(UserHomeDir(), tfchaindDir, walletsSubDir)
}

// Dir returns the directory where the data is stored for a named wallet
func Dir(name string) string {
	return filepath.Join(PersistDir(), name)
}
