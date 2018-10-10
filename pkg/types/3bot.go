package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"

	"github.com/rivine/rivine/types"
	"github.com/threefoldfoundation/tfchain/pkg/encoding"
)

const (
	// MaxNamesPerBot defines the maximum amount of names allowed per unique bot.
	MaxNamesPerBot = 5
	// MaxAddressesPerBot defines the maximum amount of addresses allowed per unique bot.
	MaxAddressesPerBot = 10
)

const (
	// BotMonth is defined as 30 days of exactly 24 hours, expressed in seconds.
	BotMonth = 60 * 60 * 24 * 30
	// MaxBotPrepaidMonths defines the amount of months that is allowed to be maximum
	// paid upfront, as to keep a 3bot active.
	MaxBotPrepaidMonths = 24
	// MaxBotPrepaidMonthsInSeconds defines the amount of time that is allowed to be maximum
	// paid upfront, which is the equavalent of roughly 2 years.
	MaxBotPrepaidMonthsInSeconds = MaxBotPrepaidMonths * BotMonth
)

var (
	// ErrTooManyBotNames is the error returned in case a bot which has more than 5
	// names defined is attempted to be (un)marshaled.
	ErrTooManyBotNames = errors.New("a 3bot can have a maximum of 5 names")
	// ErrBotNameNotUnique is the error returned in case a 3bot name is added
	// that is already registered in this 3bot.
	ErrBotNameNotUnique = errors.New("the name is already registerd with this 3bot")
	// ErrNetworkAddressNotUnique is the error returned in case a network address is added
	// that is already registered in this 3bot.
	ErrNetworkAddressNotUnique = errors.New("the network address is already registerd with this 3bot")
	// ErrBotNameDoesNotExist is the error returned in case a 3bot name is removed
	// that is not registered in this 3bot.
	ErrBotNameDoesNotExist = errors.New("the name is not registerd with this 3bot")
	// ErrNetworkAddressDoesNotExist is the error returned in case a network address is removed
	// that is not registered in this 3bot.
	ErrNetworkAddressDoesNotExist = errors.New("the network address is not registerd with this 3bot")
	// ErrTooManyBotAddresses is the error returned in case a bot which has more than 10
	// addresses defined is attempted to be (un)marshaled.
	ErrTooManyBotAddresses = errors.New("a 3bot can have a maximum of 10 addresses")
	// ErrBotExpirationExtendOverflow is returned in case a 3bot's expiration date is extended
	// using too many months (a max of 24 is allowed)
	ErrBotExpirationExtendOverflow = errors.New("a 3bot can only have up to 24 months prepaid")
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

const (
	// MaxBotID defines the maximum value a Bot ID can have,
	// in other words the biggest identifier value a Bot can have.
	MaxBotID = math.MaxUint32
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
		Addresses  []NetworkAddress `json:"addresses"`
		Names      []BotName        `json:"names"`
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
		return fmt.Errorf("BotID: %v", err)
	}
	return nil
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia,
// binary decoding this identifier from a little-endian byte slice
// of exactly 4 bytes.
func (id *BotID) UnmarshalSia(r io.Reader) error {
	x, err := encoding.UnmarshalUint32(r)
	if err != nil {
		return fmt.Errorf("BotID: %v", err)
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

// Equal returns true if this BotName and the given BotName are equal.
func (bn BotName) Equal(obn BotName) bool {
	return bytes.Compare(bn.name, obn.name) == 0
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (record BotRecord) MarshalSia(w io.Writer) error {
	enc := encoding.NewEncoder(w)

	// encode the ID and merged addr+name length
	err := enc.EncodeAll(
		record.ID,
		uint8(len(record.Addresses))|(uint8(len(record.Names))<<4),
	)
	if err != nil {
		return err
	}

	// encode all addresses and names, one after the other
	for _, addr := range record.Addresses {
		err = enc.Encode(addr)
		if err != nil {
			return err
		}
	}
	for _, name := range record.Names {
		err = enc.Encode(name)
		if err != nil {
			return err
		}
	}

	// encode the public key and the expiration date
	err = enc.EncodeAll(record.PublicKey, record.Expiration)
	if err != nil {
		return fmt.Errorf("BotRecord: MarshalSia: publicKey+expiration: %v", err)
	}
	return nil
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (record *BotRecord) UnmarshalSia(r io.Reader) error {
	decoder := encoding.NewDecoder(r)
	// decode the ID and merged addr+name len
	var pairLength uint8
	err := decoder.DecodeAll(&record.ID, &pairLength)
	if err != nil {
		return err
	}
	addrLen, nameLen := pairLength&15, pairLength>>4
	// decode all addresses
	record.Addresses = make([]NetworkAddress, 0, addrLen)
	for i := uint8(0); i < addrLen; i++ {
		var addr NetworkAddress
		err = decoder.Decode(&addr)
		if err != nil {
			return err
		}
		err = record.addNetworkAddress(addr)
		if err != nil {
			return err
		}
	}
	// decode all names
	record.Names = make([]BotName, 0, nameLen)
	for i := uint8(0); i < nameLen; i++ {
		var name BotName
		err = decoder.Decode(&name)
		if err != nil {
			return err
		}
		err = record.addName(name)
		if err != nil {
			return err
		}
	}
	err = decoder.DecodeAll(&record.PublicKey, &record.Expiration)
	if err != nil {
		return err
	}
	return nil
}

// AddNames adds one or multiple unique (DNS) names to this 3bot record.
func (record *BotRecord) AddNames(names ...BotName) error {
	var err error
	for _, name := range names {
		err = record.addName(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (record *BotRecord) addName(name BotName) error {
	// ensure the name is not already registered within this bot
	for _, rn := range record.Names {
		if rn.Equal(name) {
			return ErrBotNameNotUnique
		}
	}

	// append the name,
	// the name is not yet in this 3bot's set of names
	record.Names = append(record.Names, name)
	return nil
}

// RemoveNames removes one or multiple unique (DNS) names from this 3bot record.
func (record *BotRecord) RemoveNames(names ...BotName) error {
removeNames:
	for _, name := range names {
		for idx, rn := range record.Names {
			if rn.Equal(name) {
				record.Names = append(record.Names[:idx], record.Names[idx+1:]...)
				continue removeNames
			}
		}
		return ErrBotNameDoesNotExist
	}
	return nil
}

// AddNetworkAddresses adds one or multiple unique network addresses to this 3bot record.
func (record *BotRecord) AddNetworkAddresses(addresses ...NetworkAddress) error {
	var err error
	for _, addr := range addresses {
		err = record.addNetworkAddress(addr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (record *BotRecord) addNetworkAddress(addr NetworkAddress) error {
	// ensure the address is not already registered within this bot
	for _, raddr := range record.Addresses {
		if raddr.Equal(addr) {
			return ErrNetworkAddressNotUnique
		}
	}

	// append the address,
	// the address is not yet used in this 3bot's set of address
	record.Addresses = append(record.Addresses, addr)
	return nil
}

// RemoveNetworkAddresses removes one or multiple unique network addresses from this 3bot record.
func (record *BotRecord) RemoveNetworkAddresses(addresses ...NetworkAddress) error {
removeAddresses:
	for _, addr := range addresses {
		for idx, raddr := range record.Addresses {
			if raddr.Equal(addr) {
				record.Addresses = append(record.Addresses[:idx], record.Addresses[idx+1:]...)
				continue removeAddresses
			}
		}
		return ErrNetworkAddressDoesNotExist
	}
	return nil
}

// ExtendExpirationDate extends the expiration day of this 3bot record based on the block time
// and the months to add.
func (record *BotRecord) ExtendExpirationDate(blockTime types.Timestamp, addedMonths uint8) error {
	if addedMonths == 0 {
		return errors.New("at least one month is required in order to extend a bot's expiration date")
	}
	if addedMonths > MaxBotPrepaidMonths {
		return ErrBotExpirationExtendOverflow
	}
	bts := SiaTimestampAsCompactTimestamp(blockTime)
	if record.Expiration < bts {
		// set the block time as the base time if the record's last recorded timestamp is in the past
		record.Expiration = bts
	} // otherwise we extend based on the current graph
	newExpirationDate := record.Expiration + BotMonth*CompactTimestamp(addedMonths)
	if newExpirationDate-bts > MaxBotPrepaidMonthsInSeconds {
		return ErrBotExpirationExtendOverflow
	}
	record.Expiration = newExpirationDate
	return nil
}
