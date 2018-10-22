# Binary Encoding

Threefold Chain (tfchain) uses for most use cases the Rivine Encoding Library.
For the 3Bot transactions (and storage of 3Bot records) we use however an encoding library specific
to tfchain, to binary encode. That is what this document about.
You can read more about the Rivine encoding, also used as inspiration for the tfchain encoding library,
at <https://github.com/rivine/rivine/blob/master/doc/Encoding.md>

The main goal of the tfchain encoding library is to achieve the smallest byte footprint for encoded content.

## Standard Encoding

All integers are little-endian, encoded as unsigned integers, but the amount of types depend on the exact integral type:

| byte size | types |
| - | - |
| 1 | uint8, int8 |
| 2 | uint16, int16 |
| 3 | uint24<sup>(1)</sup> |
| 4 | uint32, int32 |
| 8 | uint64, int64, uint, int |

> (1) `uint24` is not a standard type, but the tfchain encoding lib does allow to encode uint32 integers that fit in 3 bytes, as 3 bytes.

Booleans are encoded as a single byte, `0x00` for `False` and `0x01` for `True`.

Nil pointers are equivalent to "False", i.e. a single zero byte. Valid pointers are represented by a "True" byte (0x01) followed by the encoding of the dereferenced value.

Variable-length types, such as strings and slices, are represented by a length prefix followed by the encoded value. Strings are encoded as their literal UTF-8 bytes. Slices are encoded as the concatenation of their encoded elements. The length prefix can be one, two, three or four bytes:

| byte size | inclusive size range |
| - | - |
| 1 | 0 - 127 |
| 2 | 128 - 16 383 |
| 3 | 16 384 - 2 097 151 |
| 4 | 2 097 152 - 536 870 911 |

This implies that variable-length types cannot have a size greater than `536 870 911`,
which to be fair is a very big limit for blockchain purposes. Perhaps too big of a limit already,
as it is expected that for most purposes the slice length will fit in a single byte, and the extreme cases in 2 bytes.

Maps are not supported; attempting to encode a map will cause Marshal to panic. This is because their elements are not ordered in a consistent way, and it is imperative that this encoding scheme be deterministic. To encode a map, either convert it to a slice of structs, or define a MarshalSia method (see below).

Arrays and structs are simply the concatenation of their encoded elements (no length prefix is required here as the size is fixed). Byte slices are not subject to the 8-byte integer rule; they are encoded as their literal representation, one byte per byte.

All struct fields must be exported. The ordering of struct fields is determined by their type definition.

Finally, if a type implements the SiaMarshaler interface, its MarshalSia method will be used to encode the type. Similarly, if a type implements the SiaUnmarshal interface, its UnmarshalSia method will be used to decode the type. Note that unless a type implements both interfaces, it must conform to the spec above. Otherwise, it may encode and decode itself however desired. This may be an attractive option where speed is critical, since it allows for more compact representations, and bypasses the use of reflection.

## 3Bot Encoding

3Bot transactions and records are mostly following standard encoding, but do have special encoding rules for specifc types and grouped properties. In this chapter we'll describe the types and tricks used to achieve an even compactor format, possible because of the limits put on its properties.

#### Tricks

This sub chapter explains certain common tricks used to achieve a compactor format than would be possible by just applying (the already compact) [standard encoding](#standard-encoding).

### Two slices in one

The first common trick used in 3Bot transactions is to store two slices together, by combining the length of both slices in a single prefix byte. This is possible because these slices have a maximum length that fits in 4 bits (< 16), and thus, one byte is sufficient.

### Every bit counts

One integral value that is, by definition, really small in 3Bot transactions is the "number of months". It can have a maximum of 24, meaning it fits in 5 bits (< 32), giving a waste of 3 bits.

These 3 bits are used in 3Bot transactions as flags, 1 bit per flag. The flag can indicate if certain properties are avaialble, such that it can save a byte for 0-length variable-length types (`0x00`) or a byte that would normally be used to indicate a nil-pointer (`0x00`).

### Types

#### Compact Timestamp

Timestamps are encoded in 3 bytes, this using the unix epoch time `1515000000` as the null point and recording the unix epoch time using minutes as the unit, instead of seconds.

Allowing times to be recorded up to Saturday, November 27, 2049.
Good enough for many years to come.

#### Public Key

PublicKey was already supported in Rivine, but there a 16-byte (char-array) constant is used as specifier, identifying the algorithm. In tfchain a 1-byte specifier is used instead, allowing up to 255 algorithms to be used, which is a lot more than the currently only supported algorithm.

#### Network Address

Network addresses are encoded using their raw byte representation with a 1 byte prefix. The prefix indicates the network address type as well as the length. There are 3 possibilities:

| type | name | fixed length or inclusive range |
| - | - | - |
| 0 | hostname | 0 - 63 |
| 1 | IPv4 | 4 |
| 2 | IPv6 | 16 |
| 3 | undefined |

Hostnames are enoded as raw UTF-8 encoded byte slices.
