package types

import (
	"fmt"
	"io"

	"github.com/threefoldtech/rivine/crypto"
	"github.com/threefoldtech/rivine/types"
)

// PublicKeySignaturePair pairs a public key and a signature that can be validated with it.
type PublicKeySignaturePair struct {
	PublicKey types.PublicKey `json:"publickey"`
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
	case types.SignatureAlgoEd25519:
		pksp.Signature = make(types.ByteSlice, crypto.SignatureSize)
	default:
		return fmt.Errorf("unknown SignatureAlgoType %d", pksp.PublicKey.Algorithm)
	}
	// read byte slice
	_, err = io.ReadFull(r, pksp.Signature[:])
	return err
}

// MarshalRivine implements RivineMarshaler.MarshalRivine
func (pksp PublicKeySignaturePair) MarshalRivine(w io.Writer) error {
	err := pksp.PublicKey.MarshalRivine(w)
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

// UnmarshalRivine implements RivineUnmarshaler.UnmarshalRivine
func (pksp *PublicKeySignaturePair) UnmarshalRivine(r io.Reader) error {
	// decode the public key first, which includes the algorithm type, required to know
	// what length of byte slice to expect for the signature
	err := pksp.PublicKey.UnmarshalRivine(r)
	if err != nil {
		return err
	}
	// create the expected sized byte slice, depending on the algorithm type
	switch pksp.PublicKey.Algorithm {
	case types.SignatureAlgoEd25519:
		pksp.Signature = make(types.ByteSlice, crypto.SignatureSize)
	default:
		return fmt.Errorf("unknown SignatureAlgoType %d", pksp.PublicKey.Algorithm)
	}
	// read byte slice
	_, err = io.ReadFull(r, pksp.Signature[:])
	return err
}
