package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"strconv"

	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
	"github.com/threefoldtech/rivine/types"
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
	// names defined is attempted to be (un)marshaled, or in case an amount of names
	// to be added to the bot's record would overflow this limit of 5.
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
	// addresses defined is attempted to be (un)marshaled, or in case an amount of addresses
	// to be added to the bot's record would overflow this limit of 10.
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
	// MinBotID defines the minimum value a botID can have,
	// in other words the smallest identifier value a Bot can have.
	MinBotID = 1
	// MaxBotID defines the maximum value a Bot ID can have,
	// in other words the biggest identifier value a Bot can have.
	MaxBotID = math.MaxUint32
)

type (
	// BotRecord is the record type used to store a unique 3bot in the TransactionDB.
	// Per 3bot there is one BotRecord. Once a record is created it is never deleted,
	// but it can be modified by the 3bot using one of the available Transaction types.
	BotRecord struct {
		ID         BotID                   `json:"id"`
		Addresses  NetworkAddressSortedSet `json:"addresses,omitempty"`
		Names      BotNameSortedSet        `json:"names,omitempty"`
		PublicKey  PublicKey               `json:"publickey"`
		Expiration CompactTimestamp        `json:"expiration"`
	}
)

// MarshalSia implements SiaMarshaler.MarshalSia,
// alias of MarshalRivine for backwards-compatibility reasons.
func (record BotRecord) MarshalSia(w io.Writer) error {
	return record.MarshalRivine(w)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia,
// alias of UnmarshalRivine for backwards-compatibility reasons.
func (record *BotRecord) UnmarshalSia(r io.Reader) error {
	return record.UnmarshalRivine(r)
}

// MarshalRivine implements RivineMarshaler.MarshalRivine
func (record BotRecord) MarshalRivine(w io.Writer) error {
	enc := rivbin.NewEncoder(w)

	// encode the ID and merged addr+name length
	err := enc.EncodeAll(
		record.ID,
		uint8(record.Addresses.Len())|(uint8(record.Names.Len())<<4),
	)
	if err != nil {
		return err
	}

	// encode all addresses and names, one after the other
	_, err = record.Addresses.BinaryEncode(w)
	if err != nil {
		return err
	}
	_, err = record.Names.BinaryEncode(w)
	if err != nil {
		return err
	}

	// encode the public key and the expiration date
	err = enc.EncodeAll(record.PublicKey, record.Expiration)
	if err != nil {
		return fmt.Errorf("BotRecord: MarshalRivine: publicKey+expiration: %v", err)
	}
	return nil
}

// UnmarshalRivine implements RivineUnmarshaler.UnmarshalRivine
func (record *BotRecord) UnmarshalRivine(r io.Reader) error {
	decoder := rivbin.NewDecoder(r)
	// decode the ID and merged addr+name len
	var pairLength uint8
	err := decoder.DecodeAll(&record.ID, &pairLength)
	if err != nil {
		return err
	}
	addrLen, nameLen := pairLength&15, pairLength>>4
	// decode all addresses
	err = record.Addresses.BinaryDecode(r, int(addrLen))
	if err != nil {
		return err
	}

	// decode all names
	err = record.Names.BinaryDecode(r, int(nameLen))
	if err != nil {
		return err
	}

	// decode the remaining properties
	err = decoder.DecodeAll(&record.PublicKey, &record.Expiration)
	if err != nil {
		return err
	}
	return nil
}

// AddNames adds one or multiple unique (DNS) names to this 3bot record.
func (record *BotRecord) AddNames(names ...BotName) error {
	if record.Names.Len()+len(names) > MaxNamesPerBot {
		return ErrTooManyBotNames
	}
	var err error
	for _, name := range names {
		err = record.Names.AddName(name)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveNames removes one or multiple unique (DNS) names from this 3bot record.
func (record *BotRecord) RemoveNames(names ...BotName) error {
	var err error
	for _, name := range names {
		err = record.Names.RemoveName(name)
		if err != nil {
			return err
		}
	}
	return nil
}

// ResetNames removes all names (if any) from the current record,
// reseting it to a nil set of bot names.
func (record *BotRecord) ResetNames() {
	record.Names = BotNameSortedSet{}
}

// AddNetworkAddresses adds one or multiple unique network addresses to this 3bot record.
func (record *BotRecord) AddNetworkAddresses(addresses ...NetworkAddress) error {
	if record.Addresses.Len()+len(addresses) > MaxAddressesPerBot {
		return ErrTooManyBotAddresses
	}
	var err error
	for _, addr := range addresses {
		err = record.Addresses.AddAddress(addr)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveNetworkAddresses removes one or multiple unique network addresses from this 3bot record.
func (record *BotRecord) RemoveNetworkAddresses(addresses ...NetworkAddress) error {
	var err error
	for _, addr := range addresses {
		err = record.Addresses.RemoveAddress(addr)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsExpired returns if this record indicate the bot is expired.
func (record *BotRecord) IsExpired(blockTime types.Timestamp) bool {
	return record.Expiration.SiaTimestamp() <= blockTime
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
	if record.IsExpired(blockTime) {
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

type (
	// BotID defines the identifier type for 3bots,
	// each 3bot has a unique identifier using this type.
	BotID uint32
)

type (
	// BotName defines the name type for 3bots.
	// Each 3bot can define up to 5 unique (DNS) names.
	BotName struct {
		name []byte
	}
)

// LoadString loads a botID from a string
func (id *BotID) LoadString(str string) error {
	x, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return fmt.Errorf("BotID: %v", err)
	}
	if x < MinBotID {
		return fmt.Errorf("botID has to be at least %d", MinBotID)
	}
	*id = BotID(x)
	return nil
}

// String implements fmt.Stringer.String
func (id BotID) String() string {
	return strconv.FormatUint(uint64(id), 10)
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
	return BotName{name: bytes.ToLower([]byte(name))}, nil
}

// MarshalSia implements SiaMarshaler.MarshalSia
func (bn BotName) MarshalSia(w io.Writer) error {
	return siabin.NewEncoder(w).Encode(bn.name)
}

// UnmarshalSia implements SiaUnmarshaler.UnmarshalSia
func (bn *BotName) UnmarshalSia(r io.Reader) error {
	err := siabin.NewDecoder(r).Decode(&bn.name)
	if err != nil {
		return err
	}
	bn.name = bytes.ToLower(bn.name)
	return nil
}

// MarshalRivine implements RivineMarshaler.MarshalRivine
func (bn BotName) MarshalRivine(w io.Writer) error {
	return rivbin.NewEncoder(w).Encode(bn.name)
}

// UnmarshalRivine implements RivineUnmarshaler.UnmarshalRivine
func (bn *BotName) UnmarshalRivine(r io.Reader) error {
	err := rivbin.NewDecoder(r).Decode(&bn.name)
	if err != nil {
		return err
	}
	bn.name = bytes.ToLower(bn.name)
	return nil
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

// Equals returns true if this BotName and the given BotName are equal (case insensitive).
func (bn BotName) Equals(obn BotName) bool {
	return bn.Compare(obn) == 0
}

// Compare returns an integer comparing two bot names lexicographically (case insensitive).
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func (bn BotName) Compare(obn BotName) int {
	return bytes.Compare(bn.name, obn.name)
}

type (
	// BotNameSortedSet represents a sorted set of (unique) bot names.
	//
	// A BotNameSortedSet does not expose it elements, as this is not a feature-requirement of tfchain,
	// all it aims for is to ensure the set consists only of unique elements.
	BotNameSortedSet struct {
		slice botNameSlice
	}
	botNameSlice []BotName
)

// Len returns the amount of network addresses in this sorted set.
func (bnss BotNameSortedSet) Len() int {
	return bnss.slice.Len()
}

// AddName adds a new (unique) bot name to this sorted set of bot names,
// returning an error if the name already exists within this sorted set.
func (bnss *BotNameSortedSet) AddName(name BotName) error {
	// binary search through our slice,
	// and if not found return the index where to insert the name as well
	limit := bnss.slice.Len()
	index := sort.Search(limit, func(i int) bool {
		return bnss.slice[i].Compare(name) >= 0
	})
	if index < limit && bnss.slice[index].Equals(name) {
		return ErrNetworkAddressNotUnique
	}
	// insert the new network name in the correct place
	bnss.slice = append(bnss.slice, BotName{})
	copy(bnss.slice[index+1:], bnss.slice[index:])
	bnss.slice[index] = name
	return nil
}

// RemoveName removes an existing bot name from this sorted set of bot names,
// returning an error if the name did not yet exist in this sorted set.
func (bnss *BotNameSortedSet) RemoveName(name BotName) error {
	limit := bnss.slice.Len()
	index := sort.Search(limit, func(i int) bool {
		return bnss.slice[i].Compare(name) >= 0
	})
	if index >= limit || !bnss.slice[index].Equals(name) {
		return ErrBotNameDoesNotExist
	}
	copy(bnss.slice[index:], bnss.slice[index+1:])
	bnss.slice[bnss.slice.Len()-1] = BotName{}
	bnss.slice = bnss.slice[:bnss.slice.Len()-1]
	return nil
}

// MarshalJSON implements encoding/json.Marshaler.MarshalJSON
func (bnss BotNameSortedSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(bnss.slice)
}

// UnmarshalJSON implements encoding/json.Unmarshaler.UnmarshalJSON
func (bnss *BotNameSortedSet) UnmarshalJSON(data []byte) error {
	// decode the slice
	var slice botNameSlice
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return err
	}
	// allocate suffecient memory (and erase) our internal slice
	bnss.slice = make(botNameSlice, 0, len(slice))
	// add the elements on by one, guaranteeing the names are in order and unique
	for _, addr := range slice {
		err = bnss.AddName(addr)
		if err != nil {
			return fmt.Errorf("error while unmarshaling name %v: %v", addr, err)
		}
	}
	return nil
}

// MarshalSia implements siabin.SiaMarshaler.MarshalSia
func (bnss BotNameSortedSet) MarshalSia(w io.Writer) error {
	return siabin.NewEncoder(w).Encode(bnss.slice)
}

// UnmarshalSia implements siabin.SiaUnmarshaler.UnmarshalSia
func (bnss *BotNameSortedSet) UnmarshalSia(r io.Reader) error {
	// decode the slice
	var slice botNameSlice
	err := siabin.NewDecoder(r).Decode(&slice)
	if err != nil {
		return err
	}
	// allocate suffecient memory (and erase) our internal slice
	bnss.slice = make(botNameSlice, 0, len(slice))
	// add the elements on by one, guaranteeing the names are in order and unique
	for _, addr := range slice {
		err = bnss.AddName(addr)
		if err != nil {
			return fmt.Errorf("error while unmarshaling name %v: %v", addr, err)
		}
	}
	return nil
}

// MarshalRivine implements rivbin.RivineMarshaler.MarshalRivine
func (bnss BotNameSortedSet) MarshalRivine(w io.Writer) error {
	return rivbin.NewEncoder(w).Encode(bnss.slice)
}

// UnmarshalRivine implements rivbin.RivineUnmarshaler.UnmarshalRivine
func (bnss *BotNameSortedSet) UnmarshalRivine(r io.Reader) error {
	// decode the slice
	var slice botNameSlice
	err := rivbin.NewDecoder(r).Decode(&slice)
	if err != nil {
		return err
	}
	// allocate suffecient memory (and erase) our internal slice
	bnss.slice = make(botNameSlice, 0, len(slice))
	// add the elements on by one, guaranteeing the names are in order and unique
	for _, addr := range slice {
		err = bnss.AddName(addr)
		if err != nil {
			return fmt.Errorf("error while unmarshaling name %v: %v", addr, err)
		}
	}
	return nil
}

// BinaryEncode can be used instead of MarshalRivine, should one want to
// encode the length prefix in a way other than the standard tfchain-slice approach.
// The encoding of the length has to happen prior to calling this method.
func (bnss BotNameSortedSet) BinaryEncode(w io.Writer) (int, error) {
	var (
		err     error
		encoder = rivbin.NewEncoder(w)
	)
	for _, addr := range bnss.slice {
		err = encoder.Encode(addr)
		if err != nil {
			return -1, err
		}
	}
	return bnss.slice.Len(), nil
}

// BinaryDecode can be used instead of UnmarshalRivine, should one need to
// decode the length prefix in a way other than the standard tfchain-slice approach.
// The decoding of the length has to happen prior to calling this method.
func (bnss *BotNameSortedSet) BinaryDecode(r io.Reader, length int) error {
	var (
		err     error
		decoder = rivbin.NewDecoder(r)
	)
	// allocate suffecient memory (and erase) our internal slice
	bnss.slice = make(botNameSlice, 0, length)
	// add the elements on by one, guaranteeing the names are in order and unique
	for i := 0; i < length; i++ {
		var name BotName
		err = decoder.Decode(&name)
		if err != nil {
			return err
		}
		err = bnss.AddName(name)
		if err != nil {
			return fmt.Errorf("error while unmarshaling name %v: %v", name, err)
		}
	}
	return nil
}

// Difference returns the difference of this and the other set,
// meaning it will return all bot names which are in this set but not in the other.
func (bnss BotNameSortedSet) Difference(other BotNameSortedSet) []BotName {
	indices := sortedSetDifference(bnss.Len(), other.Len(), func(a, b int) int {
		return bnss.slice[a].Compare(other.slice[b])
	})
	names := make([]BotName, 0, len(indices))
	for _, idx := range indices {
		names = append(names, bnss.slice[idx])
	}
	return names
}

// Intersection returns the intersection of this and the other set,
// meaning it will return all bot names which are in this set AND in the other.
func (bnss BotNameSortedSet) Intersection(other BotNameSortedSet) []BotName {
	indices := sortedSetIntersection(bnss.Len(), other.Len(), func(a, b int) int {
		return bnss.slice[a].Compare(other.slice[b])
	})
	names := make([]BotName, 0, len(indices))
	for _, idx := range indices {
		names = append(names, bnss.slice[idx])
	}
	return names
}

// Len implements sort.Interface.Len
func (slice botNameSlice) Len() int {
	return len(slice)
}

// Less implements sort.Interface.Less
func (slice botNameSlice) Less(i, j int) bool {
	return slice[i].Compare(slice[j]) == -1
}

// Swap implements sort.Interface.Swap
func (slice botNameSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
