package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// DenomKeyPrefix is the prefix to retrieve all Denom
	DenomKeyPrefix = "Denom/value/"
)

// DenomKey returns the store key to retrieve a Denom from the index fields
func DenomKey(
	denom string,
) []byte {
	var key []byte

	denomBytes := []byte(denom)
	key = append(key, denomBytes...)
	key = append(key, []byte("/")...)

	return key
}
