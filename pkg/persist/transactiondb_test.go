package persist

import (
	"reflect"
	"testing"
	"time"

	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	rivinetypes "github.com/threefoldtech/rivine/types"

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
		err := siabin.Unmarshal(siabin.Marshal(testCase), &result)
		if err != nil {
			t.Error(idx, "Unmarshal", err)
			continue
		}
		if !reflect.DeepEqual(testCase, result) {
			t.Error(idx, testCase, "!=", result)
		}
	}
}
func TestImplicitBotRecordUpdate_RivineMarshaling(t *testing.T) {
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
		err := rivbin.Unmarshal(rivbin.Marshal(testCase), &result)
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
