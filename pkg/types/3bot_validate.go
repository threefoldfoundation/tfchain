package types

import (
	"fmt"

	"github.com/threefoldtech/rivine/types"
)

func validateUniquenessOfNetworkAddresses(addresses []NetworkAddress) error {
	dm := make(map[string]struct{}, len(addresses))
	var (
		str    string
		exists bool
	)
	for _, addr := range addresses {
		str = addr.String()
		if _, exists = dm[str]; exists {
			return fmt.Errorf("address %s is not unique within the given slice", str)
		}
		dm[str] = struct{}{}
	}
	return nil
}

func validateUniquenessOfBotNames(names []BotName) error {
	dm := make(map[string]struct{}, len(names))
	var (
		str    string
		exists bool
	)
	for _, name := range names {
		str = name.String()
		if _, exists = dm[str]; exists {
			return fmt.Errorf("name %s is not unique within the given slice", str)
		}
		dm[str] = struct{}{}
	}
	return nil
}

func validateBotSignature(t types.Transaction, publicKey types.PublicKey, signature types.ByteSlice, ctx types.ValidationContext, extraObjects ...interface{}) error {
	condition := types.NewCondition(types.NewUnlockHashCondition(types.NewPubKeyUnlockHash(publicKey)))
	// and a matching single-signature fulfillment
	fulfillment := types.NewFulfillment(&types.SingleSignatureFulfillment{
		PublicKey: publicKey,
		Signature: signature,
	})
	// validate the signature is correct
	return condition.Fulfill(fulfillment, types.FulfillContext{
		ExtraObjects: extraObjects,
		BlockHeight:  ctx.BlockHeight,
		BlockTime:    ctx.BlockTime,
		Transaction:  t,
	})
}
