package types

import (
	"fmt"
	"strings"
)

// SpecVersion defines the version of the spec in the payload.
type SpecVersion uint64

const (
	// Unknown is an unknown spec version.
	Unknown SpecVersion = iota
	// SpecVersionBedrock is the spec version equivalent to the bellatrix release of the l1 beacon chain.
	SpecVersionBedrock
	// SpecVersionCanyon is the spec version equivalent to the capella release of the l1 beacon chain.
	SpecVersionCanyon
	// SpecVersionEcotone is the spec version equivalent to the deneb release of the l1 beacon chain.
	SpecVersionEcotone
)

var specVersionStrings = [...]string{
	"unknown",
	"bedrock",
	"canyon",
	"ecotone",
}

// MarshalJSON implements json.Marshaler.
func (d *SpecVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", specVersionStrings[*d])), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *SpecVersion) UnmarshalJSON(input []byte) error {
	var err error
	switch strings.ToLower(string(input)) {
	case `"bedrock"`:
		*d = SpecVersionBedrock
	case `"canyon"`:
		*d = SpecVersionCanyon
	case `"ecotone"`:
		*d = SpecVersionEcotone
	default:
		err = fmt.Errorf("unrecognised spec version %s", string(input))
	}

	return err
}

// String returns a string representation of the struct.
func (d SpecVersion) String() string {
	if int(d) >= len(specVersionStrings) {
		return "unknown"
	}

	return specVersionStrings[d]
}
