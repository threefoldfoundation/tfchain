package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/threefoldfoundation/tfchain/pkg/encoding"
)

const (
	// MaxNamesPerBot defines the maximum amount of names allowed per unique bot.
	MaxNamesPerBot = 5
	// MaxAddressesPerBot defines the maximum amount of addresses allowed per unique bot.
	MaxAddressesPerBot = 10
)

var (
	// ErrTooManyBotNames is the error returned in case a bot which has more than 5
	// names defined is attempted to be (un)marshaled.
	ErrTooManyBotNames = errors.New("a 3bot can have a maximum of 5 names")
	// ErrTooManyBotAddresses is the error returned in case a bot which has more than 10
	// addresses defined is attempted to be (un)marshaled.
	ErrTooManyBotAddresses = errors.New("a 3bot can have a maximum of 10 addresses")
)

const (
	// RegexpBotName is used to validate a (raw) 3bot name (string).
	RegexpBotName = `^[A-Za-z]{1}[A-Za-z\-0-9]{3,61}[A-Za-z0-9]{1}(\.[A-Za-z]{1}[A-Za-z\-0-9]{3,55}[A-Za-z0-9]{1})*$`
	// MaxLengthBotName defines the maximum length a 3bot name can have.
	MaxLengthBotName = 63
)

var (
	rexBotName = regexp.MustCompile(RegexpBotName)
)

var (
	// ErrNilBotName is the error returned in case a new bot name is attempted to be
	// created (from memory or bytes) from nil.
	ErrNilBotName = errors.New("nil bot name")
	// ErrBotNameTooLong is the error returned in case a new bot name is attempted to be
	// created using a too long string.
	ErrBotNameTooLong = errors.New("the length of a hostname can maximum be 127 bytes long")
	// ErrInvalidBotName is the error returned in case a to-be-created (or decoded)
	// botname.
	ErrInvalidBotName = errors.New("invalid bot name")
)

type (
	// BotID defines the identifier type for 3bots,
	// each 3bot has a unique identifier using this type.
	BotID uint32

	// BotName defines the name type for 3bots.
	// Each 3bot can define up to 5 unique (DNS) names.
	BotName struct {
		name []byte
	}

	// BotRecord is the record type used to store a unique 3bot in the TransactionDB.
	// Per 3bot there is one BotRecord. Once a record is created it is never deleted,
	// but it can be modified by the 3bot using one of the available Transaction types.
	BotRecord struct {
		ID         BotID            `json:"id"`
		Names      []BotName        `json:"names"`
		Addresses  []NetworkAddress `json:"addresses"`
		PublicKey  PublicKey        `json:"publickey"`
		Expiration CompactTimestamp `json:"expiration"`
	}
)

// MarshalSia implements SiaMarshaler.MarshalSia,
// binary encoding this identifier as little-endian byte
// version of its underlying uint32 representation.
func (id BotID) MarshalSia(w io.Writer) error {
	err := encoding.MarshalUint32(w, uint32(id))
	if err != nil {
		fmt.Errorf("BotID: %v", err)
	}
	return nil
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia,
// binary decoding this identifier from a little-endian byte slice
// of exactly 4 bytes.
func (id *BotID) UnmarshalSia(r io.Reader) error {
	x, err := encoding.UnmarshalUint32(r)
	if err != nil {
		fmt.Errorf("BotID: %v", err)
	}
	*id = BotID(x)
	return nil
}

// NewBotName creates a new BotName from a given (valid) string.
func NewBotName(name string) (BotName, error) {
	if name == "" {
		return BotName{}, ErrNilBotName
	}
	if len(name) > MaxLengthBotName {
		return BotName{}, ErrBotNameTooLong
	}
	if !rexBotName.MatchString(name) {
		return BotName{}, ErrInvalidBotName
	}
	return BotName{name: []byte(name)}, nil
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (bn BotName) MarshalSia(w io.Writer) error {
	return encoding.NewEncoder(w).Encode(bn.name)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (bn *BotName) UnmarshalSia(r io.Reader) error {
	return encoding.NewDecoder(r).Decode(&bn.name)
}

// String returns this BotName in a (human-readable) string format.
func (bn BotName) String() string {
	return string(bn.name)
}

// LoadString loads the BotName from a human-readable string.
func (bn *BotName) LoadString(str string) (err error) {
	*bn, err = NewBotName(str)
	return
}

// MarshalJSON marshals a byte slice as a hex string.
func (bn BotName) MarshalJSON() ([]byte, error) {
	return json.Marshal(bn.String())
}

// UnmarshalJSON decodes the json (hex-encoded) string of the byte slice.
func (bn *BotName) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	return bn.LoadString(str)
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (record BotRecord) MarshalSia(w io.Writer) error {
	// validate the length of the addresses and names, prior to starting the encoding process
	nameLen := len(record.Names)
	if nameLen > MaxNamesPerBot {
		return ErrTooManyBotNames
	}
	addrLen := len(record.Addresses)
	if addrLen > MaxAddressesPerBot {
		return ErrTooManyBotAddresses
	}
	pairLength := uint8(nameLen) | (uint8(addrLen) << 4)

	enc := encoding.NewEncoder(w)

	err := enc.Encode(record.ID)
	if err != nil {
		return fmt.Errorf("BotRecord: MarshalSia: id: %v", err)
	}

	// encode the next 2 pairs as a custom pair of tiny slices
	err = enc.Encode(pairLength)
	if err != nil {
		return fmt.Errorf("BotRecord: MarshalSia: pairLength: %v", err)
	}
	for _, name := range record.Names {
		err = enc.Encode(name)
		if err != nil {
			return fmt.Errorf("BotRecord: MarshalSia: name: %v", err)
		}
	}
	for _, addr := range record.Addresses {
		err = enc.Encode(addr)
		if err != nil {
			return fmt.Errorf("BotRecord: MarshalSia: addr: %v", err)
		}
	}

	err = enc.EncodeAll(record.PublicKey, record.Expiration)
	if err != nil {
		return fmt.Errorf("BotRecord: MarshalSia: publicKey+expiration: %v", err)
	}
	return nil
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (record *BotRecord) UnmarshalSia(r io.Reader) error {
	decoder := encoding.NewDecoder(r)
	var pairLength uint8
	err := decoder.DecodeAll(&record.ID, &pairLength)
	if err != nil {
		return fmt.Errorf("BotRecord: UnmarshalSia: id+pairLength: %v", err)
	}
	nameLen, addrLen := pairLength&15, pairLength>>4
	if nameLen > MaxNamesPerBot {
		return fmt.Errorf("BotRecord: UnmarshalSia: %v", ErrTooManyBotNames)
	}
	if addrLen > MaxAddressesPerBot {
		return fmt.Errorf("BotRecord: UnmarshalSia: %v", ErrTooManyBotAddresses)
	}
	record.Names = make([]BotName, nameLen)
	for idx := range record.Names {
		err = decoder.Decode(&record.Names[idx])
		if err != nil {
			return fmt.Errorf("BotRecord: UnmarshalSia: name: %v", err)
		}
	}
	record.Addresses = make([]NetworkAddress, addrLen)
	for idx := range record.Addresses {
		err = decoder.Decode(&record.Addresses[idx])
		if err != nil {
			return fmt.Errorf("BotRecord: UnmarshalSia: addr: %v", err)
		}
	}
	err = decoder.DecodeAll(&record.PublicKey, &record.Expiration)
	if err != nil {
		return fmt.Errorf("BotRecord: UnmarshalSia: publicKey+expiration: %v", err)
	}
	return nil
}
