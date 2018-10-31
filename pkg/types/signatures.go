package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rivine/rivine/crypto"
	"github.com/rivine/rivine/types"
	"github.com/threefoldfoundation/tfchain/pkg/encoding"
)

// SignatureAlgoType identifies a signature algorithm as a single byte.
type SignatureAlgoType uint8

const (
	// SignatureAlgoNil identifies a nil SignatureAlgoType value.
	SignatureAlgoNil SignatureAlgoType = iota
	// SignatureAlgoEd25519 identifies the Ed25519 signature Algorithm,
	// the default (and only) algorithm supported by this chain.
	SignatureAlgoEd25519
)

func (sat SignatureAlgoType) String() string {
	switch sat {
	case SignatureAlgoEd25519:
		return types.SignatureEd25519.String()
	default:
		return ""
	}
}

// LoadString loads the stringified algo type as its single byte representation.
func (sat *SignatureAlgoType) LoadString(str string) error {
	switch str {
	case types.SignatureEd25519.String():
		*sat = SignatureAlgoEd25519
	case "":
		*sat = SignatureAlgoNil
	default:
		return fmt.Errorf("unknown SignatureAlgoType string: %s", str)
	}
	return nil
}

// FromSiaPublicKey creates a PublicKey from a SiaPublicKey
func FromSiaPublicKey(spk types.SiaPublicKey) (PublicKey, error) {
	var sat SignatureAlgoType
	err := sat.LoadString(spk.Algorithm.String())
	if err != nil {
		return PublicKey{}, err
	}
	return PublicKey{Algorithm: sat, Key: spk.Key}, nil
}

// PublicKey is a public key prefixed by a Specifier. The Specifier
// indicates the algorithm used for signing and verification.
type PublicKey struct {
	Algorithm SignatureAlgoType
	Key       types.ByteSlice
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (pk PublicKey) MarshalSia(w io.Writer) error {
	err := encoding.NewEncoder(w).Encode(pk.Algorithm)
	if err != nil || pk.Algorithm == SignatureAlgoNil {
		return err // nil if pk.Algorithm == SignatureAlgoNil
	}
	l, err := w.Write([]byte(pk.Key))
	if err != nil {
		return err
	}
	if l != len(pk.Key) {
		return io.ErrShortWrite
	}
	return nil
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (pk *PublicKey) UnmarshalSia(r io.Reader) error {
	// decode the algorithm type, required to know
	// what length of byte slice to expect
	err := encoding.NewDecoder(r).Decode(&pk.Algorithm)
	if err != nil {
		return err
	}
	// create the expected sized byte slice, depending on the algorithm type
	switch pk.Algorithm {
	case SignatureAlgoEd25519:
		pk.Key = make(types.ByteSlice, crypto.PublicKeySize)
	case SignatureAlgoNil:
		pk.Key = nil
	default:
		return fmt.Errorf("unknown SignatureAlgoType %d", pk.Algorithm)
	}
	// read byte slice
	_, err = io.ReadFull(r, pk.Key[:])
	return err
}

// LoadString is the inverse of SiaPublicKey.String().
func (pk *PublicKey) LoadString(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return errors.New("invalid public key string")
	}
	err := pk.Key.LoadString(parts[1])
	if err != nil {
		return err
	}
	return pk.Algorithm.LoadString(parts[0])
}

// String defines how to print a PublicKey.
func (pk PublicKey) String() string {
	return pk.Algorithm.String() + ":" + pk.Key.String()
}

// MarshalJSON marshals a byte slice as a hex string.
func (pk PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pk.String())
}

// UnmarshalJSON decodes the json (hex-encoded) string of the byte slice.
func (pk *PublicKey) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	return pk.LoadString(str)
}

// SiaPublicKey returns this PublicKey as a SiaPublicKey
func (pk PublicKey) SiaPublicKey() (types.SiaPublicKey, error) {
	switch pk.Algorithm {
	case SignatureAlgoEd25519:
		return types.SiaPublicKey{
			Algorithm: types.SignatureEd25519,
			Key:       pk.Key,
		}, nil
	default:
		return types.SiaPublicKey{}, fmt.Errorf("unknown algorithm type: %d", pk.Algorithm)
	}
}

// PublicKeySignaturePair pairs a public key and a signature that can be validated with it.
type PublicKeySignaturePair struct {
	PublicKey PublicKey       `json:"publickey"`
	Signature types.ByteSlice `json:"signature"`
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (pksp PublicKeySignaturePair) MarshalSia(w io.Writer) error {
	err := pksp.PublicKey.MarshalSia(w)
	if err != nil {
		return err
	}
	l, err := w.Write([]byte(pksp.Signature))
	if err != nil {
		return err
	}
	if l != len(pksp.Signature) {
		return io.ErrShortWrite
	}
	return nil
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (pksp *PublicKeySignaturePair) UnmarshalSia(r io.Reader) error {
	// decode the public key first, which includes the algorithm type, required to know
	// what length of byte slice to expect for the signature
	err := pksp.PublicKey.UnmarshalSia(r)
	if err != nil {
		return err
	}
	// create the expected sized byte slice, depending on the algorithm type
	switch pksp.PublicKey.Algorithm {
	case SignatureAlgoEd25519:
		pksp.Signature = make(types.ByteSlice, crypto.SignatureSize)
	default:
		return fmt.Errorf("unknown SignatureAlgoType %d", pksp.PublicKey.Algorithm)
	}
	// read byte slice
	_, err = io.ReadFull(r, pksp.Signature[:])
	return err
}
