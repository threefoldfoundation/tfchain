# Binary Encoding

Threefold Chain (tfchain) uses for most use cases the Sia (Binary) Encoding Library.
For the 3Bot transactions (and storage of 3Bot records) we use however an encoding library specific
to rivine, to binary encode. That is what this document about.
You can read more about the Sia Binary encoding
at <https://github.com/threefoldtech/rivine/blob/master/doc/encoding/SiaEncoding.md>.
You can read more about the Rivine Binary encoding
at <https://github.com/threefoldtech/rivine/blob/master/doc/encoding/RivineEncoding.md>,
used for the 3Bot transactions and inspired by the Sia Binary Encoding.

The main goal of the rivine encoding library is to achieve the smallest byte footprint for encoded content.

## 3Bot Encoding

3Bot transactions and records are mostly following standard (rivine) encoding, but do have special encoding rules for specific types and grouped properties. In this chapter we'll describe the types and tricks used to achieve an even compactor format, possible because of the limits put on its properties.

#### Tricks

This sub chapter explains certain common tricks used to achieve a compactor format than would be possible by just applying (the already compact) [standard encoding](#standard-encoding).

### Two slices in one

The first common trick used in 3Bot transactions is to store two slices together, by combining the length of both slices in a single prefix byte. This is possible because these slices have a maximum length that fits in 4 bits (< 16), and thus, one byte is sufficient.

### Every bit counts

One integral value that is, by definition, really small in 3Bot transactions is the "number of months". It can have a maximum of 24, meaning it fits in 5 bits (< 32), giving a waste of 3 bits.

These 3 bits are used in 3Bot transactions as flags, 1 bit per flag. The flag can indicate if certain properties are available, such that it can save a byte for 0-length variable-length types (`0x00`) or a byte that would normally be used to indicate a nil-pointer (`0x00`).

### Types

#### Compact Timestamp

Timestamps are encoded in 3 bytes, this using the unix epoch time `1515000000` as the null point and recording the unix epoch time using minutes as the unit, instead of seconds.

Allowing times to be recorded up to Saturday, November 27, 2049.
Good enough for many years to come.

#### Public Key

PublicKey was already supported in Rivine, but there a 16-byte (character array) constant is used as specifier, identifying the algorithm. In tfchain a 1-byte specifier is used instead, allowing up to 255 algorithms to be used, which is a lot more than the currently only supported algorithm.

#### Network Address

Network addresses are encoded using their raw byte representation with a 1 byte prefix. The prefix indicates the network address type as well as the length. There are 3 possibilities:

| type | name | fixed length or inclusive range |
| - | - | - |
| 0 | hostname | 0 - 63 |
| 1 | IPv4 | 4 |
| 2 | IPv6 | 16 |
| 3 | undefined |

Hostnames are encoded as raw UTF-8 encoded byte slices.
