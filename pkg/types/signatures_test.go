package types

import (
	"encoding/json"
	"testing"

	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
)

const psTestPairJSON = `{
	"publickey": "ed25519:1a91300a1cca5faab7f3967a3dde84b7d2c6cfb96a60221b37eb743fd4765588",
	"signature": "1468358ed24f812535f15964bbb93c02c5c14fdc5809da7dac83a510249dc702d51076c795150966cf16f361f0c7b53c529813a569ed57cc2e0e691783f3660d"
}`

func TestPublicKeySignaturePairMarshalRivine(t *testing.T) {
	var pksp PublicKeySignaturePair
	err := json.Unmarshal([]byte(psTestPairJSON), &pksp)
	if err != nil {
		t.Fatal(err)
	}

	b := rivbin.Marshal(pksp)

	var pksp2 PublicKeySignaturePair
	err = rivbin.Unmarshal(b, &pksp2)
	if err != nil {
		t.Fatal(err)
	}

	b, err = json.MarshalIndent(pksp2, "", "	")
	if err != nil {
		t.Fatal(err)
	}

	out := string(b)
	if out != psTestPairJSON {
		t.Fatal("unexpected JSON output:", out, "!=", psTestPairJSON)
	}
}

func TestPublicKeySignaturePairMarshalSia(t *testing.T) {
	var pksp PublicKeySignaturePair
	err := json.Unmarshal([]byte(psTestPairJSON), &pksp)
	if err != nil {
		t.Fatal(err)
	}

	b := siabin.Marshal(pksp)

	var pksp2 PublicKeySignaturePair
	err = siabin.Unmarshal(b, &pksp2)
	if err != nil {
		t.Fatal(err)
	}

	b, err = json.MarshalIndent(pksp2, "", "	")
	if err != nil {
		t.Fatal(err)
	}

	out := string(b)
	if out != psTestPairJSON {
		t.Fatal("unexpected JSON output:", out, "!=", psTestPairJSON)
	}
}
