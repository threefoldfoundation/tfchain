package main

import (
	"fmt"
	"testing"

	"github.com/rivine/rivine/types"
)

func TestMultiSignatureConditionIsStandardCondition(t *testing.T) {
	ctx := types.StandardCheckContext{}
	// create the condition manually
	msc := MultiSignatureCondition{
		MultiSignatureCondition: &types.MultiSignatureCondition{
			UnlockHashes: []types.UnlockHash{
				unlockHashFromHex("01746677df456546d93729066dd88514e2009930f3eebac3c93d43c88a108f8f9aa9e7c6f58893"),
				unlockHashFromHex("01c46a8e1e7f1bb0e3b7ec6c93b9c4f3e5d89e855f5a57f22d478d72d6233391153fac7d179087"),
			},
			MinimumSignatureCount: 1,
		},
	}
	// ensure that the internal condition's standard check does pass
	err := msc.MultiSignatureCondition.IsStandardCondition(ctx)
	if err != nil {
		t.Fatal("expected standard condition check pass, but it failed: ", err)
	}
	// ensure that our internal condition's standard check fails
	err = msc.IsStandardCondition(ctx)
	if err == nil {
		t.Fatal("expected standard condition check to fail, but it didn't")
	}
}

func TestRegisteredMultiSignatureCondition(t *testing.T) {
	// temporary overwrite multisig condition type, just for this unit test
	RegisteredBlockHeightLimitedMultiSignatureCondition()
	defer types.RegisterUnlockConditionType(types.ConditionTypeMultiSignature,
		func() types.MarshalableUnlockCondition { return new(types.MultiSignatureCondition) })

	const jsonCondition = `{
	"type": 4,
	"data": {
		"unlockhashes": [
			"01746677df456546d93729066dd88514e2009930f3eebac3c93d43c88a108f8f9aa9e7c6f58893",
			"01c46a8e1e7f1bb0e3b7ec6c93b9c4f3e5d89e855f5a57f22d478d72d6233391153fac7d179087"
		],
		"minimumsignaturecount": 2
	}
}`

	// decode our json-encoded multisig condition
	var condition types.UnlockConditionProxy
	err := condition.UnmarshalJSON([]byte(jsonCondition))
	if err != nil {
		t.Fatal("failed to decode multisignature condition into proxy condition: ", err)
	}

	// ensure the condition type is MultiSig
	if ct := condition.ConditionType(); ct != types.ConditionTypeMultiSignature {
		t.Fatalf("expected condition type to be %d, but it was %d instead",
			types.ConditionTypeMultiSignature, ct)
	}

	// sanity check, ensure it is our type
	if _, ok := condition.Condition.(*MultiSignatureCondition); !ok {
		t.Fatalf("expected condition type to be (our) *MultiSignatureCondition, but it was %T instead",
			condition.Condition)
	}

	// ensure that it can't be used yet at height 0
	ctx := types.StandardCheckContext{}
	err = condition.IsStandardCondition(ctx)
	if err == nil {
		t.Fatal("expected standard condition check to fail, but it didn't")
	}

	// ensure that it can be used at the minimum height
	ctx.BlockHeight = MinimumBlockHeightForMultiSignatureConditions
	err = condition.IsStandardCondition(ctx)
	if err != nil {
		t.Fatal("expected standard condition check pass, but it failed: ", err)
	}
}

func unlockHashFromHex(hstr string) (uh types.UnlockHash) {
	err := uh.LoadString(hstr)
	if err != nil {
		panic(fmt.Sprintf("func unlockHashFromHex(%s) failed: %v", hstr, err))
	}
	return
}
