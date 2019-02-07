package main

import (
	"fmt"
	"strconv"
	"strings"

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
	amount, err := cc.ParseCoinString(args[1])
	if err != nil {
		return err
	}

	var to types.UnlockHash
	err = to.LoadString(args[0])
	if err != nil {
		return err
	}

	w, err := wallet.Load(walletName, backend)
	if err != nil {
		return err
	}

	err = w.TransferCoins(amount, to, []byte(cmds.DataString), cmds.GenerateNewRefundAddress)
	if err != nil {
		return err
	}

	fmt.Println("Transaction posted")
	fmt.Println("Transfered", cc.ToCoinStringWithUnit(amount), "to", to.String())
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
