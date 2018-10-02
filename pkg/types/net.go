package types

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"regexp"

	"github.com/threefoldfoundation/tfchain/pkg/encoding"
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
	ErrHostnameTooLong = errors.New("the length of a hostname can maximum be 127 bytes long")
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
func (na NetworkAddress) MarshalSia(w io.Writer) error {
	length := len(na.addr)
	err := encoding.MarshalUint8(w, uint8(na.t)|uint8(length)<<2)
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

// UnmarshalSia unmarshals this NetworkAddress from a compact binary format.
func (na *NetworkAddress) UnmarshalSia(r io.Reader) error {
	lengthAndType, err := encoding.UnmarshalUint8(r)
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
	if na.t == NetworkAddressHostname {
		return string(na.addr) // also covers a nil address
	}
	// given we assume our NetworkAddress is valid, the type must be IPv4 or IPv6
	return net.IP(na.addr).String()
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
