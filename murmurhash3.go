// MurmurHash3 was written by Austin Appleby, and is placed in the public
package murmurhash3

import (
	"encoding/binary"
	"hash"
)

type (
	murmurhash3A uint32
	murmurhash3C uint32
	murmurhash3F uint64
)

func New3A() hash.Hash32 {
	var m murmurhash3A = 0
	return &m
}
func (m *murmurhash3A) Reset() { *m = 0 }
func (m *murmurhash3A) Size() int {
	return 4
}
func (m *murmurhash3A) Write(p []byte) (n int, err error) {
	*m = murmurhash3A(Murmur3A(p, uint32(*m)))
	return len(p), nil
}
func (m *murmurhash3A) Sum32() uint32 {
	return uint32(*m)
}
func (m *murmurhash3A) Sum() []byte {
	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, uint32(*m))
	return p
}

func New3C() hash.Hash32 {
	var m murmurhash3C = 0
	return &m
}
func (m *murmurhash3C) Reset() { *m = 0 }
func (m *murmurhash3C) Size() int {
	return 4
}
func (m *murmurhash3C) Write(p []byte) (n int, err error) {
	*m = murmurhash3C(Murmur3C(p, uint32(*m))[0])
	return len(p), nil
}
func (m *murmurhash3C) Sum32() uint32 {
	return uint32(*m)
}
func (m *murmurhash3C) Sum() []byte {
	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, uint32(*m))
	return p
}

func New3F() hash.Hash64 {
	var m murmurhash3F = 0
	return &m
}
func (m *murmurhash3F) Reset() { *m = 0 }
func (m *murmurhash3F) Size() int {
	return 8
}
func (m *murmurhash3F) Write(p []byte) (n int, err error) {
	*m = murmurhash3F(Murmur3F(p, uint64(*m))[0])
	return len(p), nil
}
func (m *murmurhash3F) Sum64() uint64 {
	return uint64(*m)
}
func (m *murmurhash3F) Sum() []byte {
	p := make([]byte, 8)
	binary.BigEndian.PutUint64(p, uint64(*m))
	return p
}

func rotl32(x uint32, r uint8) uint32 {
	return (x << r) | (x >> (32 - r))
}

func rotl64(x uint64, r uint8) uint64 {
	return (x << r) | (x >> (64 - r))
}

func fmix32(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}

func fmix64(h uint64) uint64 {
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	h *= 0xc4ceb9fe1a85ec53
	h ^= h >> 33

	return h
}

func Murmur3A(key []byte, seed uint32) uint32 {
	nblocks := len(key) / 4
	var h1 uint32 = seed

	var c1 uint32 = 0xcc9e2d51
	var c2 uint32 = 0x1b873593

	// body
	for i := 0; i < nblocks; i++ {
		k1 := binary.LittleEndian.Uint32(key[i*4:]) // TODO Validate

		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2

		h1 ^= k1
		h1 = rotl32(h1, 13)
		h1 = h1*5 + 0xe6546b64
	}

	// tail
	var k1 uint32 = 0
	var tail []byte = key[nblocks*4:] // TODO Validate
	switch len(key) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2
		h1 ^= k1
	}

	//finalization
	h1 ^= uint32(len(key))

	h1 = fmix32(h1)

	return h1
}

func Murmur3C(key []byte, seed uint32) [4]uint32 {
	nblocks := len(key) / 16
	var h1 uint32 = seed
	var h2 uint32 = seed
	var h3 uint32 = seed
	var h4 uint32 = seed

	var c1 uint32 = 0x239b961b
	var c2 uint32 = 0xab0e9789
	var c3 uint32 = 0x38b34ae5
	var c4 uint32 = 0xa1e38b93

	// body
	for i := 0; i < nblocks; i++ {
		k1 := binary.LittleEndian.Uint32(key[(i*4+0)*4:]) // TODO Validate
		k2 := binary.LittleEndian.Uint32(key[(i*4+1)*4:])
		k3 := binary.LittleEndian.Uint32(key[(i*4+2)*4:])
		k4 := binary.LittleEndian.Uint32(key[(i*4+3)*4:])

		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2
		h1 ^= k1

		h1 = rotl32(h1, 19)
		h1 += h2
		h1 = h1*5 + 0x561ccd1b

		k2 *= c2
		k2 = rotl32(k2, 16)
		k2 *= c3
		h2 ^= k2

		h2 = rotl32(h2, 17)
		h2 += h3
		h2 = h2*5 + 0x0bcaa747

		k3 *= c3
		k3 = rotl32(k3, 17)
		k3 *= c4
		h3 ^= k3

		h3 = rotl32(h3, 15)
		h3 += h4
		h3 = h3*5 + 0x96cd1c35

		k4 *= c4
		k4 = rotl32(k4, 18)
		k4 *= c1
		h4 ^= k4

		h4 = rotl32(h4, 13)
		h4 += h1
		h4 = h4*5 + 0x32ac3b17
	}

	// tail
	var k1 uint32 = 0
	var k2 uint32 = 0
	var k3 uint32 = 0
	var k4 uint32 = 0
	var tail []byte = key[nblocks*16:] // TODO Validate
	switch len(key) & 15 {
	case 15:
		k4 ^= uint32(tail[14]) << 16
		fallthrough
	case 14:
		k4 ^= uint32(tail[13]) << 8
		fallthrough
	case 13:
		k4 ^= uint32(tail[12]) << 0
		k4 *= c4
		k4 = rotl32(k4, 18)
		k4 *= c1
		h4 ^= k4
		fallthrough
	case 12:
		k3 ^= uint32(tail[11]) << 24
		fallthrough
	case 11:
		k3 ^= uint32(tail[10]) << 16
		fallthrough
	case 10:
		k3 ^= uint32(tail[9]) << 8
		fallthrough
	case 9:
		k3 ^= uint32(tail[8]) << 0
		k3 *= c3
		k3 = rotl32(k3, 17)
		k3 *= c4
		h3 ^= k3
		fallthrough
	case 8:
		k2 ^= uint32(tail[7]) << 24
		fallthrough
	case 7:
		k2 ^= uint32(tail[6]) << 16
		fallthrough
	case 6:
		k2 ^= uint32(tail[5]) << 8
		fallthrough
	case 5:
		k2 ^= uint32(tail[4]) << 0
		k2 *= c2
		k2 = rotl32(k2, 16)
		k2 *= c3
		h2 ^= k2
		fallthrough
	case 4:
		k1 ^= uint32(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0]) << 0
		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2
		h1 ^= k1
	}

	//finalization
	h1 ^= uint32(len(key))
	h2 ^= uint32(len(key))
	h3 ^= uint32(len(key))
	h4 ^= uint32(len(key))

	h1 += h2
	h1 += h3
	h1 += h4
	h2 += h1
	h3 += h1
	h4 += h1

	h1 = fmix32(h1)
	h2 = fmix32(h2)
	h3 = fmix32(h3)
	h4 = fmix32(h4)

	h1 += h2
	h1 += h3
	h1 += h4
	h2 += h1
	h3 += h1
	h4 += h1

	return [4]uint32{h1, h2, h3, h4}
}

func Murmur3F(key []byte, seed uint64) [2]uint64 {
	nblocks := len(key) / 16
	var h1 uint64 = seed
	var h2 uint64 = seed

	var c1 uint64 = 0x87c37b91114253d5
	var c2 uint64 = 0x4cf5ad432745937f

	// body
	for i := 0; i < nblocks; i++ {
		k1 := binary.LittleEndian.Uint64(key[(i*2+0)*8:]) // TODO Validate
		k2 := binary.LittleEndian.Uint64(key[(i*2+1)*8:])

		k1 *= c1
		k1 = rotl64(k1, 31)
		k1 *= c2
		h1 ^= k1

		h1 = rotl64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2
		k2 = rotl64(k2, 33)
		k2 *= c1
		h2 ^= k2

		h2 = rotl64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	// tail
	var k1 uint64 = 0
	var k2 uint64 = 0
	var tail []byte = key[nblocks*16:] // TODO Validate
	switch len(key) & 15 {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8]) << 0
		k2 *= c2
		k2 = rotl64(k2, 33)
		k2 *= c1
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0]) << 0
		k1 *= c1
		k1 = rotl64(k1, 31)
		k1 *= c2
		h1 ^= k1
	}

	//finalization
	h1 ^= uint64(len(key))
	h2 ^= uint64(len(key))

	h1 += h2
	h2 += h1

	h1 = fmix64(h1)
	h2 = fmix64(h2)

	h1 += h2
	h2 += h1

	return [2]uint64{h1, h2}
}
