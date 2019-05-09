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

type (
	// ReservationData is the structure for the data to include in a transaction
	// used to create a reservation
	ReservationData struct {
		Email    string   `json:"email"`
		Size     int      `json:"size"`
		Type     Workload `json:"type"`
		Location string
	}

	// Workload is the shorthand identifier for a workload
	Workload string
)

// Constants representing the different kind of workloads
const (
	VM Workload = "vm"
	S3          = "s3"
)

const (
	// PublicBroker is the address as string of the public broker
	PublicBroker = "019bc85e0d710d928f163cbe9bf9f4911462488468ab66b758e178ea7ef978992fc203130127b7"
)

func (cmds *cmds) walletInit(cmd *cobra.Command, args []string) error {
	wallet, err := wallet.New(args[0], cmds.KeysToLoad, cmds.Network)
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
	wallet, err := wallet.NewWalletFromMnemonic(args[0], strings.TrimSpace(mnemonic), cmds.KeysToLoad, cmds.Network)
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
	w, err := wallet.Load(walletName)
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

	w, err := wallet.Load(walletName)
	if err != nil {
		return err
	}

	cts, err := w.GetChainConstants()
	if err != nil {
		return err
	}

	cc := client.NewCurrencyConvertor(types.CurrencyUnits{OneCoin: cts.OneCoin}, cts.ChainInfo.CoinUnit)
	amount, err := cc.ParseCoinString(args[0])
	if err != nil {
		return err
	}
	amounts := []types.Currency{amount}

	var targetConditions []types.MarshalableUnlockCondition
	if len(args) == 2 {
		// nil condition
		if args[1] == "" {
			targetConditions = []types.MarshalableUnlockCondition{&types.NilCondition{}}
		} else {
			// actual address
			var to types.UnlockHash
			err = to.LoadString(args[1])
			if err != nil {
				return err
			}
			targetConditions = []types.MarshalableUnlockCondition{types.NewUnlockHashCondition(to)}
		}
	} else {
		if cmds.MultiSig {
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
			targetConditions = []types.MarshalableUnlockCondition{types.NewMultiSignatureCondition(addresses, sigAmt)}
		} else {
			if len(args)%2 != 0 {
				return errors.New("amount/address pair(s) expected")
			}
			amounts = []types.Currency{}
			for i := 0; i < len(args); i += 2 {
				amount, err := cc.ParseCoinString(args[i])
				if err != nil {
					return err
				}
				amounts = append(amounts, amount)
				addr := types.UnlockHash{}
				err = addr.LoadString(args[i+1])
				if err != nil {
					return err
				}
				targetConditions = append(targetConditions, types.NewUnlockHashCondition(addr))
			}
		}
	}

	if cmds.LockString != "" {
		timeLock, err := parseLockTime(cmds.LockString)
		if err != nil {
			return err
		}
		for i := range targetConditions {
			targetConditions[i] = types.NewTimeLockCondition(timeLock, targetConditions[i])
		}
	}

	targetConditionProxies := make([]types.UnlockConditionProxy, 0, len(targetConditions))
	addresses := make([]string, 0, len(targetConditions))
	for i, condition := range targetConditions {
		targetConditionProxies = append(targetConditionProxies, types.NewCondition(condition))
		addresses = append(addresses, targetConditionProxies[i].UnlockHash().String())
	}

	err = w.TransferCoinsMulti(amounts, targetConditionProxies, []byte(cmds.DataString), cmds.GenerateNewRefundAddress)
	if err != nil {
		return err
	}

	fmt.Println("Transaction posted")
	for i, address := range addresses {
		amount := amounts[i]
		fmt.Println("Transfered", cc.ToCoinStringWithUnit(amount), "to", address)
	}
	return nil
}

func (cmds *cmds) walletReserveVM(cmd *cobra.Command, args []string) error {
	// a nodeid has a length of 12
	if len(args[1]) != 12 {
		return errors.New("Invalid node ID length")
	}

	walletName := cmd.Parent().Parent().Name()
	w, err := wallet.Load(walletName)
	if err != nil {
		return err
	}

	return reserveWorkload(w, VM, args[0], args[1], args[2], cmds.Broker, cmds.GenerateNewRefundAddress)
}

func (cmds *cmds) walletReserveS3(cmd *cobra.Command, args []string) error {
	walletName := cmd.Parent().Parent().Name()
	w, err := wallet.Load(walletName)
	if err != nil {
		return err
	}

	return reserveWorkload(w, S3, args[0], args[1], args[2], cmds.Broker, cmds.GenerateNewRefundAddress)
}

func reserveWorkload(w *wallet.Wallet, workload Workload, sizeString string,
	location string, email string, customBroker string, newRefundAddr bool) error {

	// use the user defined broker if set
	broker := PublicBroker
	if customBroker != "" {
		broker = customBroker
	}

	cts, err := w.GetChainConstants()
	if err != nil {
		return err
	}
	cc := client.NewCurrencyConvertor(types.CurrencyUnits{OneCoin: cts.OneCoin}, cts.ChainInfo.CoinUnit)

	var to types.UnlockHash
	err = to.LoadString(broker)
	if err != nil {
		return err
	}
	targetCondition := types.NewUnlockHashCondition(to)

	size, cost, err := parseTypeSize(workload, sizeString)
	if err != nil {
		return err
	}
	amount, _ := cc.ParseCoinString(strconv.Itoa(cost))
	data := ReservationData{
		Email:    email,
		Type:     workload,
		Size:     size,
		Location: location,
	}

	buf, err := encodeReservationData(data)
	if err != nil {
		return err
	}

	err = w.TransferCoins(amount, types.NewCondition(targetCondition), buf, newRefundAddr)
	if err != nil {
		return err
	}

	fmt.Println("Reservation created")
	fmt.Printf("Paid %v to %v to reserve a %v of size %v\n", cc.ToCoinStringWithUnit(amount),
		targetCondition.UnlockHash().String(), workload, size)

	return nil
}

func (cmds *cmds) walletAddresses(cmd *cobra.Command, args []string) error {
	walletName := cmd.Parent().Name()
	w, err := wallet.Load(walletName)
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

	w, err := wallet.Load(walletName)
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

	w, err := wallet.Load(walletName)
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

func parseTypeSize(typ Workload, sizeString string) (int, int, error) {
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		return 0, 0, err
	}

	switch typ {
	case VM:
		switch size {
		case 1:
			return size, 1, nil
		case 2:
			return size, 4, nil
		default:
			return 0, 0, fmt.Errorf("Invalid size %v for 'vm', only size '1' and '2' supported", size)
		}
	case S3:
		switch size {
		case 1:
			return size, 10, nil
		case 2:
			return size, 40, nil
		default:
			return 0, 0, fmt.Errorf("Invalid size %v for 's3', only size '1' and '2' supported", size)
		}
	default:
		return 0, 0, fmt.Errorf("Invalid type '%v', only 'vm' and 's3' are supported", typ)
	}
}

// encodeReservationData converts a ReservationData struct to a byteslice
// which can be included in a transaction
// data layout:
// 1 byte type
// 1 byte size
// 1 byte lenght of nodeID or farm name, depending on type
// nodeID for VM, farm name for S3
// 1 byte length of email address
// email address
func encodeReservationData(data ReservationData) ([]byte, error) {
	bytes := make([]byte, 2)

	if data.Type == VM {
		bytes[0] = 1
	} else if data.Type == S3 {
		bytes[0] = 2
	} else {
		return nil, fmt.Errorf("Unknown workload %s", data.Type)
	}

	bytes[1] = byte(data.Size)

	bytes = append(bytes, byte(len(data.Location)))
	bytes = append(bytes, []byte(data.Location)...)

	bytes = append(bytes, byte(len(data.Email)))
	bytes = append(bytes, []byte(data.Email)...)

	return bytes, nil

}
