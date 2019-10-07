package eth

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/log"
)

//InitializeKeystore imports, creates  or opens the keystore and unlock the first account
//If no account is found in the keystore, one is created.
func InitializeKeystore(datadir string, accountJSON string, accountPass string) (*keystore.KeyStore, error) {
	ks := keystore.NewKeyStore(filepath.Join(datadir, "keys"), keystore.StandardScryptN, keystore.StandardScryptP)
	var acc accounts.Account
	if accountJSON != "" {
		log.Info("Importing account")
		blob, err := ioutil.ReadFile(accountJSON)
		if err != nil {
			return nil, err
		}
		// Import the account
		acc, err = ks.Import(blob, accountPass, accountPass)
		if err != nil {
			return nil, err
		}

	} else {
		if len(ks.Accounts()) == 0 {
			log.Info("Creating a new account")
			var err error
			acc, err = ks.NewAccount(accountPass)
			if err != nil {
				return nil, err
			}

		} else {
			log.Info("Loading existing account")
			acc = ks.Accounts()[0]
		}
	}
	log.Info("Unlocking account " + acc.Address.String())
	if err := ks.Unlock(acc, accountPass); err != nil {
		return nil, err
	}
	log.Info("Account unlocked")
	return ks, nil
}
