package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/threefoldtech/rivine/pkg/client"
	"github.com/threefoldtech/rivine/types"

	"github.com/bgentry/speakeasy"

	"github.com/spf13/cobra"
	"github.com/threefoldfoundation/tfchain/cmd/light-client/explorer"
	"github.com/threefoldfoundation/tfchain/cmd/light-client/wallet"
)

func (cmds *cmds) walletInit(cmd *cobra.Command, args []string) error {
	wallet, err := wallet.New(args[0], cmds.KeysToLoad)
	if err != nil {
		return err
	}
	mnemonic, err := wallet.Mnemonic()
	if err != nil {
		return err
	}
	fmt.Println("Created new wallet", args[0], "!")
	fmt.Println("Wallet seed:")
	fmt.Println(mnemonic)
	return nil
}

func (cmds *cmds) walletRecover(cmd *cobra.Command, args []string) error {
	mnemonic, err := speakeasy.Ask("Seed:")
	if err != nil {
		return err
	}
	wallet, err := wallet.NewWalletFromMnemonic(args[0], strings.TrimSpace(mnemonic), cmds.KeysToLoad)
	if err != nil {
		return err
	}
	newmnemonic, err := wallet.Mnemonic()
	if err != nil {
		return err
	}
	if newmnemonic != mnemonic {
		panic("Different mnemonics")
	}
	fmt.Println("Created wallet", args[0], "from existing seed")
	fmt.Println("Wallet seed:")
	fmt.Println(newmnemonic)
	return nil
}

func (cmds *cmds) walletSeed(cmd *cobra.Command, args []string) error {
	walletName := cmd.Parent().Name()
	w, err := wallet.Load(walletName, nil)
	if err != nil {
		return err
	}
	mnemonic, err := w.Mnemonic()
	if err != nil {
		return err
	}
	fmt.Println("Wallet seed:")
	fmt.Println(mnemonic)
	return nil
}

func (cmds *cmds) walletSend(cmd *cobra.Command, args []string) error {
	walletName := cmd.Parent().Name()

	backend := explorer.NewTestnetGroupedExplorer()
	cts, err := backend.GetChainConstants()
	if err != nil {
		return err
	}

	cc := client.NewCurrencyConvertor(types.CurrencyUnits{OneCoin: cts.OneCoin}, cts.ChainInfo.CoinUnit)
	amount, err := cc.ParseCoinString(args[0])
	if err != nil {
		return err
	}

	w, err := wallet.Load(walletName, backend)
	if err != nil {
		return err
	}

	var targetCondition types.MarshalableUnlockCondition
	if len(args) == 2 {
		// nil condition
		if args[1] == "" {
			targetCondition = &types.NilCondition{}
		} else {
			// actual address
			var to types.UnlockHash
			err = to.LoadString(args[1])
			if err != nil {
				return err
			}
			targetCondition = types.NewUnlockHashCondition(to)
		}
	} else {
		addressCount := len(args) - 1
		var addresses []types.UnlockHash
		// try to parse the last argument as an amount of signatures
		sigAmt, err := parseAmount(args[len(args)-1])
		if err != nil {
			// all multisig addresses
			sigAmt = uint64(addressCount)
		} else {
			// last input is an amount so we have 1 less address input
			addressCount--
			if sigAmt > uint64(addressCount) {
				return errors.New("Invalid amount of signatures required, can't require more signatures than there are addresses")
			}
		}
		// first arg is the amount of tokens so ignore that
		for i := 1; i < addressCount+1; i++ {
			addr := types.UnlockHash{}
			err = addr.LoadString(args[i])
			if err != nil {
				return err
			}
			addresses = append(addresses, addr)
		}
		targetCondition = types.NewMultiSignatureCondition(addresses, sigAmt)
	}

	if cmds.LockString != "" {
		timeLock, err := parseLockTime(cmds.LockString)
		if err != nil {
			return err
		}
		targetCondition = types.NewTimeLockCondition(timeLock, targetCondition)
	}

	err = w.TransferCoins(amount, types.NewCondition(targetCondition), []byte(cmds.DataString), cmds.GenerateNewRefundAddress)
	if err != nil {
		return err
	}

	fmt.Println("Transaction posted")
	fmt.Println("Transfered", cc.ToCoinStringWithUnit(amount), "to", targetCondition.UnlockHash().String())
	return nil
}

func (cmds *cmds) walletAddresses(cmd *cobra.Command, args []string) error {
	walletName := cmd.Parent().Name()
	w, err := wallet.Load(walletName, nil)
	if err != nil {
		return err
	}
	addresses := w.ListAddresses()

	fmt.Println("Wallet addresses:")
	for _, addr := range addresses {
		fmt.Println(addr)
	}
	return nil
}

func (cmds *cmds) walletBalance(cmd *cobra.Command, args []string) error {
	walletName := cmd.Name()

	backend := explorer.NewTestnetGroupedExplorer()
	cts, err := backend.GetChainConstants()
	if err != nil {
		return err
	}

	cc := client.NewCurrencyConvertor(types.CurrencyUnits{OneCoin: cts.OneCoin}, cts.ChainInfo.CoinUnit)

	w, err := wallet.Load(walletName, backend)
	if err != nil {
		return err
	}

	fmt.Println("Checking wallet balance")
	fmt.Println("Depending on the amount of addresses you have loaded, this may take a while")
	fmt.Println("")

	unlockedBalance, lockedBalance, err := w.GetBalance()
	if err != nil {
		return err
	}

	fmt.Println("Wallet balance:")
	fmt.Println("Unlocked:\t", cc.ToCoinStringWithUnit(unlockedBalance))
	fmt.Println("Locked:  \t", cc.ToCoinStringWithUnit(lockedBalance))
	return nil
}

func (cmds *cmds) walletLoad(cmd *cobra.Command, args []string) error {
	walletName := cmd.Parent().Parent().Name()

	// generate 1 address if no additional arg is specified
	amountString := "1"
	if len(args) > 0 {
		amountString = args[0]
	}
	amount, err := strconv.ParseUint(amountString, 10, 64)
	if err != nil {
		return err
	}

	w, err := wallet.Load(walletName, nil)
	if err != nil {
		return err
	}

	return w.LoadKeys(amount)
}

func parseAmount(amt string) (uint64, error) {
	return strconv.ParseUint(amt, 10, 64)
}

func parseLockTime(lockStr string) (uint64, error) {
	// block height or unix time stamp
	integer, err := strconv.ParseUint(lockStr, 10, 64)
	if err == nil {
		return integer, err
	}

	// date
	timestamp, err := time.Parse("_2 Jan 2006", lockStr)
	if err == nil {
		return uint64(timestamp.Unix()), nil
	}

	// date time
	timestamp, err = time.Parse("_2 Jan 2006 15:04", lockStr)
	if err == nil {
		return uint64(timestamp.Unix()), nil
	}

	// duration
	duration, err := time.ParseDuration(lockStr)
	if err == nil {
		return uint64(time.Now().Add(duration).Unix()), nil
	}

	return 0, errors.New("Unrecognized locktime")
}
