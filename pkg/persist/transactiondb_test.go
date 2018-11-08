package persist

import (
	"reflect"
	"testing"
	"time"

	rivinetypes "github.com/threefoldtech/rivine/types"

	"github.com/threefoldfoundation/tfchain/pkg/encoding"
	"github.com/threefoldfoundation/tfchain/pkg/types"
)

func TestImplicitBotRecordUpdate_SiaMarshaling(t *testing.T) {
	testCases := []implicitBotRecordUpdate{
		{},
		{
			PreviousExpirationTime: types.CompactTimestampNullpoint,
		},
		{
			PreviousExpirationTime: types.SiaTimestampAsCompactTimestamp(rivinetypes.Timestamp(time.Now().Unix())),
		},
		{
			InactiveNamesRemoved: []types.BotName{
				mustNewBotName(t, "bbbbb.aaaaa"),
			},
		},
		{
			InactiveNamesRemoved: []types.BotName{
				mustNewBotName(t, "aaaaa.bbbbb"),
				mustNewBotName(t, "bbbbb.aaaaa"),
			},
		},
		{
			PreviousExpirationTime: types.SiaTimestampAsCompactTimestamp(rivinetypes.Timestamp(time.Now().Unix())),
			InactiveNamesRemoved: []types.BotName{
				mustNewBotName(t, "aaaaa.bbbbb"),
				mustNewBotName(t, "bbbbb.aaaaa"),
			},
		},
	}
	for idx, testCase := range testCases {
		var result implicitBotRecordUpdate
		err := encoding.Unmarshal(encoding.Marshal(testCase), &result)
		if err != nil {
			t.Error(idx, "Unmarshal", err)
			continue
		}
		if !reflect.DeepEqual(testCase, result) {
			t.Error(idx, testCase, "!=", result)
		}
	}
}

func mustNewBotName(t *testing.T, str string) types.BotName {
	t.Helper()
	name, err := types.NewBotName(str)
	if err != nil {
		t.Fatal(err)
	}
	return name
}
