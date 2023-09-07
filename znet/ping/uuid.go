// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/zqlpaopao/tool/string-byte/src"
	"io"
)

type UID [16]byte

var (
	grander = rand.Reader // random function
	NilUID  UID           // empty UUID, all zeros
	Nil     UUID          // empty UUID, all zeros
)

// NewUID creates a new random UUID or panics.  New is equivalent to
// the expression
//
//	uuid.Must(uuid.NewRandom())
func NewUID() UID {
	return NewRandomUID()
}

// NewRandomUID returns a Random (Version 4) UUID.
//
// The strength of the UUIDs is based on the strength of the crypto/rand
// package.
//
// A note about uniqueness derived from the UUID Wikipedia entry:
//
//	Randomly generated UUIDs have 122 random bits.  One's annual risk of being
//	hit by a meteorite is estimated to be one chance in 17 billion, that
//	means the probability is about 0.00000000006 (6 × 10−11),
//	equivalent to the odds of creating a few tens of trillions of UUIDs in a
//	year and having one duplicate.
func NewRandomUID() UID {
	return NewRandomUIDFromReader(grander)
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (u UID) String() string {
	return src.Bytes2String(u[:])
}

// NewRandomUIDFromReader returns a UUID based on bytes read from a given io.Reader.
func NewRandomUIDFromReader(r io.Reader) UID {
	var uuid UID
	_, err := io.ReadFull(r, uuid[:])
	if err != nil {
		return NilUID
	}
	//uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	//uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}

/********************************** UUID ************************************/

// A UUID is a 128 bit (16 byte) Universal Unique Identifier as defined in RFC
// 4122.
type UUID [8]byte

// NewUUid creates a new random UUID or panics.  New is equivalent to
// the expression
//
//	uuid.Must(uuid.NewRandom())
func NewUUid() UUID {
	return NewRandom()
}

// NewRandom returns a Random (Version 4) UUID.
//
// The strength of the UUIDs is based on the strength of the crypto/rand
// package.
//
// A note about uniqueness derived from the UUID Wikipedia entry:
//
//	Randomly generated UUIDs have 122 random bits.  One's annual risk of being
//	hit by a meteorite is estimated to be one chance in 17 billion, that
//	means the probability is about 0.00000000006 (6 × 10−11),
//	equivalent to the odds of creating a few tens of trillions of UUIDs in a
//	year and having one duplicate.
func NewRandom() UUID {
	return NewRandomFromReader(grander)
}

// NewRandomFromReader returns a UUID based on bytes read from a given io.Reader.
func NewRandomFromReader(r io.Reader) UUID {
	var uuid UUID
	_, err := io.ReadFull(r, uuid[:])
	if err != nil {
		return Nil
	}
	//uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	//uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (uuid UUID) String() string {
	var buf = make([]byte, 16)
	encodeHex(buf[:], uuid)
	return src.Bytes2String(buf[:])
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:])
	//dst[8] = '-'
	//hex.Encode(dst[9:13], uuid[4:6])
	//dst[13] = '-'
	//hex.Encode(dst[14:18], uuid[6:8])
	//dst[18] = '-'
	//hex.Encode(dst[19:23], uuid[8:10])
	//dst[23] = '-'
	//hex.Encode(dst[24:], uuid[10:])
}

// MarshalBinary implements encoding.BinaryMarshaller.
func (uuid UUID) MarshalBinary() ([]byte, error) {
	return uuid[:], nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (uuid UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("invalid UUID (got %d bytes)", len(data))
	}
	copy(uuid[:], data)
	return nil
}
