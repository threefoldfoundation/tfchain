package main

import (
	"fmt"

	"github.com/rivine/rivine/types"
)

// RegisteredBlockHeightLimitedMultiSignatureCondition registers the multisig condition,
// and thus implicitly the fulfillment as well, in a way that it is limited to a certain block height.
func RegisteredBlockHeightLimitedMultiSignatureCondition() {
	types.RegisterUnlockConditionType(types.ConditionTypeMultiSignature,
		func() types.MarshalableUnlockCondition { return new(MultiSignatureCondition) })
}

// MultiSignatureCondition wraps around the Rivine-standard MultiSignatureCondition type,
// as to ensure that in the standard network of tfchain, it can only be used since blockheight 42000
type MultiSignatureCondition struct {
	*types.MultiSignatureCondition
}

const (
	// MinimumBlockHeightForMultiSignatureConditions defines the blockheight
	// since when MultiSignatureConditions are considered standard on the tfchain standard network.
	MinimumBlockHeightForMultiSignatureConditions = types.BlockHeight(42000)
)

// IsStandardCondition implements UnlockCondition.IsStandardCondition,
// wrapping around the internal MultiSignatureCondition's IsStandardCondition check,
// adding a pre-check of the blockheight
func (msc MultiSignatureCondition) IsStandardCondition(ctx types.StandardCheckContext) error {
	if ctx.BlockHeight < MinimumBlockHeightForMultiSignatureConditions {
		return fmt.Errorf(
			"multisignature conditions are only allowed since blockheight %d",
			MinimumBlockHeightForMultiSignatureConditions)
	}
	return msc.MultiSignatureCondition.IsStandardCondition(ctx)
}
