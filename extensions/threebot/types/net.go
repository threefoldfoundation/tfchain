package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"sort"

	"github.com/threefoldtech/rivine/pkg/encoding/rivbin"
	"github.com/threefoldtech/rivine/pkg/encoding/siabin"
)

const (
	// RegexpHostname is used to validate a (raw) hostname (string).
	RegexpHostname = `^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`
	// MaxLengthHostname defines the maximum length a hostname can have,
	// within the context of tfchain.
	MaxLengthHostname = 63
)

// NetworkAddressType defines the type of a network address.
type NetworkAddressType uint8

const (
	// NetworkAddressHostname represents a valid hostname, assumed to be a valid FQDN,
	// and defined as described in RFC 1178.
	NetworkAddressHostname NetworkAddressType = iota
	// NetworkAddressIPv4 represents an IPv4 address, meaning an address identified by 4 bytes,
	// and defined as described in RFC 791.
	NetworkAddressIPv4
	// NetworkAddressIPv6 represents an IPv6 address, meaning an address identified by 6 bytes,
	// and defined as described in RFC 2460.
	NetworkAddressIPv6
)

var (
	// ErrNilHostname is the error returned in case a new network address is attempted to be
	// created (from memory or bytes) from nil.
	ErrNilHostname = errors.New("nil hostname")
	// ErrHostnameTooLong is the error returned in case a new network address is attempted to be
	// created with a too long string.
	ErrHostnameTooLong = errors.New("the length of a hostname can maximum be 63 bytes long")
	// ErrInvalidNetworkAddress is the error returned in case a to-be-created (or decoded)
	// network address is invalid (meaning it is no valid hostname or IPv4/IPv6 address).
	ErrInvalidNetworkAddress = errors.New("invalid network address")
)

var (
	rexHostname = regexp.MustCompile(RegexpHostname)
)

// NetworkAddress represents a NetworkAddress,
// meaning an IPv4/6 address or (domain) hostname.
type NetworkAddress struct {
	t    NetworkAddressType
	addr []byte
}

var (
	v4InV6Prefix    = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}
	v4InV6PrefixLen = len(v4InV6Prefix)
)

// NewNetworkAddress creates a new NetworkAddress from a given (valid) string.
func NewNetworkAddress(addr string) (NetworkAddress, error) {
	if addr == "" {
		return NetworkAddress{}, ErrNilHostname
	}
	if rexHostname.MatchString(addr) {
		if len(addr) > MaxLengthHostname {
			return NetworkAddress{}, ErrHostnameTooLong
		}
		return NetworkAddress{
			t:    NetworkAddressHostname,
			addr: []byte(addr),
		}, nil
	}
	na := NetworkAddress{addr: []byte(net.ParseIP(addr))}
	if len(na.addr) > v4InV6PrefixLen && bytes.Equal(na.addr[:v4InV6PrefixLen], v4InV6Prefix) {
		na.addr = na.addr[v4InV6PrefixLen:]
	}
	switch len(na.addr) {
	case 4:
		na.t = NetworkAddressIPv4
	case 16:
		na.t = NetworkAddressIPv6
	default:
		return NetworkAddress{}, ErrInvalidNetworkAddress
	}
	return na, nil
}

// MarshalSia marshals this NetworkAddress in a compact binary format.
// Alias of MarshalRivine, for backwards-compatibility
func (na NetworkAddress) MarshalSia(w io.Writer) error {
	return na.MarshalRivine(w)
}

// UnmarshalSia unmarshals this NetworkAddress from a semi-compact binary format.
// Alias of UnmarshalRivine, for backwards-compatibility
func (na *NetworkAddress) UnmarshalSia(r io.Reader) error {
	return na.UnmarshalRivine(r)
}

// MarshalRivine marshals this NetworkAddress in a compact binary format.
func (na NetworkAddress) MarshalRivine(w io.Writer) error {
	length := len(na.addr)
	err := rivbin.MarshalUint8(w, uint8(na.t)|uint8(length)<<2)
	if err != nil {
		return err
	}
	n, err := w.Write(na.addr)
	if err != nil {
		return err
	}
	if n != length {
		return io.ErrShortWrite
	}
	return nil
}

// UnmarshalRivine unmarshals this NetworkAddress from a compact binary format.
func (na *NetworkAddress) UnmarshalRivine(r io.Reader) error {
	lengthAndType, err := rivbin.UnmarshalUint8(r)
	if err != nil {
		return err
	}
	length := lengthAndType >> 2
	switch na.t = NetworkAddressType(lengthAndType & 3); na.t {
	case NetworkAddressHostname:
		if length > MaxLengthHostname {
			return ErrInvalidNetworkAddress
		}
	case NetworkAddressIPv4:
		if length != 4 {
			return ErrInvalidNetworkAddress
		}
	case NetworkAddressIPv6:
		if length != 16 {
			return ErrInvalidNetworkAddress
		}
	}
	na.addr = make([]byte, length)
	_, err = io.ReadFull(r, na.addr)
	if err != nil {
		return err
	}
	if na.t == NetworkAddressHostname && !rexHostname.Match(na.addr) {
		// a hostname can be invalid however,
		// given that we marshal the string directly
		return ErrInvalidNetworkAddress
	}
	return nil
}

// String returns this NetworkAddress in a (human-readable) string format.
func (na NetworkAddress) String() string {
	switch na.t {
	case NetworkAddressIPv4:
		p := make(net.IP, net.IPv6len)
		copy(p, v4InV6Prefix)
		copy(p[v4InV6PrefixLen:], na.addr)
		return p.String()
	case NetworkAddressIPv6:
		return net.IP(na.addr).String()
	default: // NetworkAddressHostname
		return string(na.addr) // also covers a nil address
	}
}

// LoadString loads the NetworkAddress from a human-readable string.
func (na *NetworkAddress) LoadString(str string) (err error) {
	*na, err = NewNetworkAddress(str)
	return
}

// MarshalJSON marshals a byte slice as a hex string.
func (na NetworkAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(na.String())
}

// UnmarshalJSON decodes the json (hex-encoded) string of the byte slice.
func (na *NetworkAddress) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	return na.LoadString(str)
}

// Equals returns true if this NetworkAddress and the given NetworkAddress are equal.
func (na NetworkAddress) Equals(ona NetworkAddress) bool {
	return na.t == ona.t && bytes.Compare(na.addr, ona.addr) == 0
}

// Compare returns an integer comparing two network addresses.
// If the types are equal the addresses are compared lexicographically,
// otherwise the compare result of the network address types is returned.
// The final result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func (na NetworkAddress) Compare(ona NetworkAddress) int {
	if na.t < ona.t {
		return -1
	}
	if na.t > ona.t {
		return 1
	}
	return bytes.Compare(na.addr, ona.addr)
}

type (
	// NetworkAddressSortedSet represents a sorted set of (unique) network addresses.
	//
	// A NetworkAddressSortedSet does not expose it elements, as this is not a feature-requirement of tfchain,
	// all it aims for is to ensure the set consists only of unique elements.
	NetworkAddressSortedSet struct {
		slice networkAddressSlice
	}
	networkAddressSlice []NetworkAddress
)

// Len returns the amount of network addresses in this sorted set.
func (nass NetworkAddressSortedSet) Len() int {
	return nass.slice.Len()
}

// AddAddress adds a new (unique) network address to this sorted set of network addresses,
// returning an error if the address already exists within this sorted set.
func (nass *NetworkAddressSortedSet) AddAddress(address NetworkAddress) error {
	// binary search through our slice,
	// and if not found return the index where to insert the address as well
	limit := nass.slice.Len()
	index := sort.Search(limit, func(i int) bool {
		return nass.slice[i].Compare(address) >= 0
	})
	if index < limit && nass.slice[index].Equals(address) {
		return ErrNetworkAddressNotUnique
	}
	// insert the new network address in the correct place
	nass.slice = append(nass.slice, NetworkAddress{})
	copy(nass.slice[index+1:], nass.slice[index:])
	nass.slice[index] = address
	return nil
}

// RemoveAddress removes an existing network address from this sorted set of network addresses,
// returning an error if the address did not yet exist in this sorted set.
func (nass *NetworkAddressSortedSet) RemoveAddress(address NetworkAddress) error {
	limit := nass.slice.Len()
	index := sort.Search(limit, func(i int) bool {
		return nass.slice[i].Compare(address) >= 0
	})
	if index >= limit || !nass.slice[index].Equals(address) {
		return ErrNetworkAddressDoesNotExist
	}
	copy(nass.slice[index:], nass.slice[index+1:])
	nass.slice[nass.slice.Len()-1] = NetworkAddress{}
	nass.slice = nass.slice[:nass.slice.Len()-1]
	return nil
}

// MarshalJSON implements encoding/json.Marshaler.MarshalJSON
func (nass NetworkAddressSortedSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(nass.slice)
}

// UnmarshalJSON implements encoding/json.Unmarshaler.UnmarshalJSON
func (nass *NetworkAddressSortedSet) UnmarshalJSON(data []byte) error {
	// decode the slice
	var slice networkAddressSlice
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return err
	}
	// allocate suffecient memory (and erase) our internal slice
	nass.slice = make(networkAddressSlice, 0, len(slice))
	// add the elements on by one, guaranteeing the addresses are in order and unique
	for _, addr := range slice {
		err = nass.AddAddress(addr)
		if err != nil {
			return fmt.Errorf("error while unmarshaling addr %v: %v", addr, err)
		}
	}
	return nil
}

// MarshalSia implements siabin.SiaMarshaler.MarshalSia
func (nass NetworkAddressSortedSet) MarshalSia(w io.Writer) error {
	return siabin.NewEncoder(w).Encode(nass.slice)
}

// UnmarshalSia implements siabin.SiaUnmarshaler.UnmarshalSia
func (nass *NetworkAddressSortedSet) UnmarshalSia(r io.Reader) error {
	// decode the slice
	var slice networkAddressSlice
	err := siabin.NewDecoder(r).Decode(&slice)
	if err != nil {
		return err
	}
	// allocate suffecient memory (and erase) our internal slice
	nass.slice = make(networkAddressSlice, 0, len(slice))
	// add the elements on by one, guaranteeing the addresses are in order and unique
	for _, addr := range slice {
		err = nass.AddAddress(addr)
		if err != nil {
			return fmt.Errorf("error while unmarshaling addr %v: %v", addr, err)
		}
	}
	return nil
}

// MarshalRivine implements rivbin.RivineMarshaler.MarshalRivine
func (nass NetworkAddressSortedSet) MarshalRivine(w io.Writer) error {
	return rivbin.NewEncoder(w).Encode(nass.slice)
}

// UnmarshalRivine implements rivbin.RivineUnmarshaler.UnmarshalRivine
func (nass *NetworkAddressSortedSet) UnmarshalRivine(r io.Reader) error {
	// decode the slice
	var slice networkAddressSlice
	err := rivbin.NewDecoder(r).Decode(&slice)
	if err != nil {
		return err
	}
	// allocate suffecient memory (and erase) our internal slice
	nass.slice = make(networkAddressSlice, 0, len(slice))
	// add the elements on by one, guaranteeing the addresses are in order and unique
	for _, addr := range slice {
		err = nass.AddAddress(addr)
		if err != nil {
			return fmt.Errorf("error while unmarshaling addr %v: %v", addr, err)
		}
	}
	return nil
}

// BinaryEncode can be used instead of MarshalRivine, should one want to
// encode the length prefix in a way other than the standard tfchain-slice approach.
// The encoding of the length has to happen prior to calling this method.
func (nass NetworkAddressSortedSet) BinaryEncode(w io.Writer) (int, error) {
	var (
		err     error
		encoder = rivbin.NewEncoder(w)
	)
	for _, addr := range nass.slice {
		err = encoder.Encode(addr)
		if err != nil {
			return -1, err
		}
	}
	return nass.slice.Len(), nil
}

// BinaryDecode can be used instead of UnmarshalRivine, should one need to
// decode the length prefix in a way other than the standard tfchain-slice approach.
// The decoding of the length has to happen prior to calling this method.
func (nass *NetworkAddressSortedSet) BinaryDecode(r io.Reader, length int) error {
	var (
		err     error
		decoder = rivbin.NewDecoder(r)
	)
	// allocate suffecient memory (and erase) our internal slice
	nass.slice = make(networkAddressSlice, 0, length)
	// add the elements on by one, guaranteeing the addresses are in order and unique
	for i := 0; i < length; i++ {
		var addr NetworkAddress
		err = decoder.Decode(&addr)
		if err != nil {
			return err
		}
		err = nass.AddAddress(addr)
		if err != nil {
			return fmt.Errorf("error while unmarshaling addr %v: %v", addr, err)
		}
	}
	return nil
}

// Len implements sort.Interface.Len
func (slice networkAddressSlice) Len() int {
	return len(slice)
}

// Less implements sort.Interface.Less
func (slice networkAddressSlice) Less(i, j int) bool {
	return slice[i].Compare(slice[j]) == -1
}

// Swap implements sort.Interface.Swap
func (slice networkAddressSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
