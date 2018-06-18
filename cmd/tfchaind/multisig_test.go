package main

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/rivine/rivine/encoding"
	"github.com/rivine/rivine/types"
)

func TestMultiSignatureConditionIsStandardCondition(t *testing.T) {
	ctx := types.StandardCheckContext{}
	// create the condition manually
	msc := MultiSignatureCondition{
		MultiSignatureCondition: types.MultiSignatureCondition{
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

func TestDecodeBinaryCoinOutputsForIssue141(t *testing.T) {
	// temporary overwrite multisig condition type, just for this unit test
	RegisteredBlockHeightLimitedMultiSignatureCondition()
	defer types.RegisterUnlockConditionType(types.ConditionTypeMultiSignature,
		func() types.MarshalableUnlockCondition { return new(types.MultiSignatureCondition) })

	const binaryHexData = "0200000000000000050000000000000009cd5b050004520000000000000002000000000000000200000000000000017115d8f27e0ff38b77766fb9838e0a7736cea38ac00ef12347fac04ba71710dc0149a5496fea27315b7db6251e5dfda23bc9d4bf677c5a5c2d70f1382c44357197060000000000000002b0aa9e4a00012100000000000000017115d8f27e0ff38b77766fb9838e0a7736cea38ac00ef12347fac04ba71710dc"
	var coinoutputs []types.CoinOutput
	binaryData, err := hex.DecodeString(binaryHexData)
	if err != nil {
		t.Fatal("failed to hex-decode binary data", err)
	}
	err = encoding.Unmarshal(binaryData, &coinoutputs)
	if err != nil {
		t.Fatal("failed to binary-decode coin outputs", err)
	}
}

func TestDecodeBinaryTransactionSetForIssue141(t *testing.T) {
	// temporary overwrite multisig condition type, just for this unit test
	RegisteredBlockHeightLimitedMultiSignatureCondition()
	defer types.RegisterUnlockConditionType(types.ConditionTypeMultiSignature,
		func() types.MarshalableUnlockCondition { return new(types.MultiSignatureCondition) })

	const binaryHexData = "01000000000000000185010000000000000100000000000000107df606f88a99943f290b54a2815dd0ca6eb051f8534444e51439f3d11455ab018000000000000000656432353531390000000000000000002000000000000000b5662caa078efd42b25f3ab10768b55fd0607ed8cb8e3c44f3b26df1d17ef93440000000000000001220697d9acae414dd60b216f6372144c66265b506b008933dd125bb7ae621bc2a476a575917ac2e82310bd0e361957fc7907af116e296020dd0837b1aefd2000200000000000000050000000000000009cd5b050004520000000000000002000000000000000200000000000000017115d8f27e0ff38b77766fb9838e0a7736cea38ac00ef12347fac04ba71710dc0149a5496fea27315b7db6251e5dfda23bc9d4bf677c5a5c2d70f1382c44357197060000000000000002b0aa9e4a00012100000000000000017115d8f27e0ff38b77766fb9838e0a7736cea38ac00ef12347fac04ba71710dc000000000000000000000000000000000100000000000000040000000000000005f5e1000000000000000000"
	var transactions []types.Transaction
	binaryData, err := hex.DecodeString(binaryHexData)
	if err != nil {
		t.Fatal("failed to hex-decode binary data", err)
	}
	err = encoding.Unmarshal(binaryData, &transactions)
	if err != nil {
		t.Fatal("failed to binary-decode transactions", err)
	}
}
