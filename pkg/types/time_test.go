package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math"
	"testing"

	"github.com/rivine/rivine/types"
	"github.com/threefoldfoundation/tfchain/pkg/encoding"
)

func TestCompactTimestampLimits(t *testing.T) {
	testCases := []CompactTimestamp{
		// lower limit
		CompactTimestampNullpoint,
		// values in between
		CompactTimestampNullpoint + 1*CompactTimestampAccuracyInSeconds,
		CompactTimestampNullpoint + 1*CompactTimestampAccuracyInSeconds + 5,
		CompactTimestampNullpoint + 42*CompactTimestampAccuracyInSeconds,
		CompactTimestampNullpoint + 100000*CompactTimestampAccuracyInSeconds,
		CompactTimestampNullpoint + 100000*CompactTimestampAccuracyInSeconds + CompactTimestampAccuracyInSeconds/2,
		CompactTimestampNullpoint + 1234321*CompactTimestampAccuracyInSeconds,
		CompactTimestampNullpoint + (math.MaxUint32>>9)*CompactTimestampAccuracyInSeconds,
		// upper limit
		CompactTimestampNullpoint + (math.MaxUint32>>8)*CompactTimestampAccuracyInSeconds,
	}
	for idx, testCase := range testCases {
		// expected value for all limit tests on this test case
		expected := testCase - (testCase % CompactTimestampAccuracyInSeconds)

		// Test SiaTimestampAsCompactTimestamp Limits
		cts := SiaTimestampAsCompactTimestamp(types.Timestamp(testCase))
		if cts != expected {
			t.Error(idx+1, "SiaTimestampAsCompactTimestamp", "unexpected unmarshal result:", cts, "!=", expected)
		}

		// Test BinaryEncoding Limits
		err := encoding.Unmarshal(encoding.Marshal(testCase), &cts)
		if err != nil {
			t.Error(idx+1, "unmarshal error", testCase, "message:", err)
			continue
		}
		if cts != expected {
			t.Error(idx+1, "encoding.Unmarshal(encoding.Marshal())", "unexpected unmarshal result:", cts, "!=", expected)
		}

		// Test JSONEncoding Limits
		b, err := json.Marshal(testCase)
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(b, &cts)
		if err != nil {
			t.Fatal(err)
		}
		if cts != expected {
			t.Error(idx+1, "json.Unmarshal(json.Marshal())", "unexpected unmarshal result:", cts, "!=", expected)
		}
	}
}

func TestCompactTimestampBinaryEncodingUnmarshalMarshalExample(t *testing.T) {
	const hexStr = `7af905`
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatal(err)
	}
	var cs CompactTimestamp
	err = cs.UnmarshalSia(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	buffer := bytes.NewBuffer(nil)
	err = cs.MarshalSia(buffer)
	if err != nil {
		t.Fatal(err)
	}
	str := hex.EncodeToString(buffer.Bytes())
	if str != hexStr {
		t.Fatal("unexpected hex result", str)
	}
}
