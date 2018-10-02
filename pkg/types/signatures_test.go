package types

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/threefoldfoundation/tfchain/pkg/encoding"

	"github.com/rivine/rivine/types"
)

func TestSignatureAlgoTypeSiaMarshaling(t *testing.T) {
	for sat := SignatureAlgoType(0); sat < 255; sat++ {
		b := encoding.Marshal(sat)
		if len(b) == 0 {
			t.Error(sat, "encoding.Marshal", "<nil>")
			continue
		}
		if len(b) != 1 {
			t.Error(sat, len(b), "!= 1")
		}
		var result SignatureAlgoType
		err := encoding.Unmarshal(b, &result)
		if err != nil {
			t.Error(sat, "encoding.Unmarshal", err)
			continue
		}
		if result != sat {
			t.Error(result, "!=", sat)
		}
	}
}

func TestSignatureAlgoTypeJSONMarshaling(t *testing.T) {
	for sat := SignatureAlgoType(0); sat < 255; sat++ {
		b, err := json.Marshal(sat)
		if err != nil || len(b) == 0 {
			t.Error(sat, "json.Marshal", err)
			continue
		}
		var result SignatureAlgoType
		err = json.Unmarshal(b, &result)
		if err != nil {
			t.Error(sat, "json.Unmarshal", err)
			continue
		}
		if result != sat {
			t.Error(result, "!=", sat)
		}
	}
}

var exampleSiaPublicKeyStrings = []string{
	"ed25519:97d784b93d5769d2df0010b793622eae14a33992b958e8406cceb827a8101d29",
	"ed25519:cb859ec8da13d0bcfc7b1c3c8e6647b5510791eda3a74ba1bba4954a8c74e4a9",
	"ed25519:857c029d8689c97d51f314e1a4e6a4543c42a696ee93ce848d1247bf24eb52a3",
	"ed25519:4683705f729a65e9e133e1719d05ad8ac45a14e44fcf6c85de19e5ac7fcd2e9d",
}

func TestPublicKeySiaBiDirectionality(t *testing.T) {
	for idx, example := range exampleSiaPublicKeyStrings {
		var spk types.SiaPublicKey
		err := spk.LoadString(example)
		if err != nil {
			t.Error(idx, "(*SiaPublicKey).LoadString", example, err)
			continue
		}
		pk, err := FromSiaPublicKey(spk)
		if err != nil {
			t.Error(idx, "FromSiaPublicKey", example, err)
			continue
		}
		spk, err = pk.SiaPublicKey()
		if err != nil {
			t.Error(idx, "SiaPublicKey", example, err)
			continue
		}
		str := spk.String()
		if str != example {
			t.Error(idx, str, "!=", example)
		}
	}
}

func TestFromSiaPublicKey(t *testing.T) {
	for idx, example := range exampleSiaPublicKeyStrings {
		var spk types.SiaPublicKey
		err := spk.LoadString(example)
		if err != nil {
			t.Error(idx, "(*SiaPublicKey).LoadString", example, err)
			continue
		}
		pk, err := FromSiaPublicKey(spk)
		if err != nil {
			t.Error(idx, "FromSiaPublicKey", example, err)
			continue
		}
		str := pk.String()
		if str != example {
			t.Error(idx, str, "!=", example)
		}
	}
}

func TestPublicKeyLoadStringString(t *testing.T) {
	for idx, example := range exampleSiaPublicKeyStrings {
		var pk PublicKey
		err := pk.LoadString(example)
		if err != nil {
			t.Error(idx, "LoadString", example, err)
			continue
		}
		str := pk.String()
		if str != example {
			t.Error(idx, str, "!=", example)
		}
	}
}

func TestPublicKeySiaBinaryMarshaling(t *testing.T) {
	for idx, example := range exampleSiaPublicKeyStrings {
		var pk PublicKey
		err := pk.LoadString(example)
		if err != nil {
			t.Error(idx, "LoadString", example, err)
			continue
		}

		// start binary marshal
		b := encoding.Marshal(pk)
		if len(b) == 0 {
			t.Error(idx, "encoding.Marshal", "<nil>")
		}
		err = encoding.Unmarshal(b, &pk)
		if err != nil {
			t.Error(idx, "encoding.Unmarshal", err)
		}
		// end binary marshal

		str := pk.String()
		if str != example {
			t.Error(idx, str, "!=", example)
		}
	}
}

func TestPublicKeySiaJSONMarshaling(t *testing.T) {
	for idx, example := range exampleSiaPublicKeyStrings {
		var pk PublicKey
		err := pk.LoadString(example)
		if err != nil {
			t.Error(idx, "LoadString", example, err)
			continue
		}

		// start JSON marshal
		b, err := json.Marshal(pk)
		if err != nil || len(b) == 0 {
			t.Error(idx, "encoding.Marshal", err)
		}
		err = json.Unmarshal(b, &pk)
		if err != nil {
			t.Error(idx, "encoding.Unmarshal", err)
		}
		// end JSON marshal

		str := pk.String()
		if str != example {
			t.Error(idx, str, "!=", example)
		}
	}
}

func TestPublicKeyEd25519BinaryEncodedByteLength(t *testing.T) {
	const str = `ed25519:97d784b93d5769d2df0010b793622eae14a33992b958e8406cceb827a8101d29`
	var pk PublicKey
	err := pk.LoadString(str)
	if err != nil {
		t.Error(err)
	}
	b := encoding.Marshal(pk)
	if len(b) != 34 {
		t.Error(len(b), "!=", 34, hex.EncodeToString(b))
	}
}
