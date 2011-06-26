package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"hash"
	"os"
)

const (
	m    uint32 = 0x5bd1e995
	r    uint32 = 24
	Size = 4
)

type MurmurHash2A struct {
	hash  uint32
	tail  uint32
	count uint32
	size  uint32
}

func mmix(h, k uint32) (uint32, uint32) {
	k *= m
	k ^= k >> r
	k *= m
	h *= m
	h ^= k
	return h, k
}

func (me *MurmurHash2A) mixTail(data []byte, i, length *int) {
	for *length != 0 && (*length < 4 || me.count != 0) {
		me.tail |= uint32(data[*i]) << (me.count * 8)

		me.count++
		*i++
		*length--

		if me.count == 4 {
			me.hash, me.tail = mmix(me.hash, me.tail)
			me.tail = 0
			me.count = 0
		}
	}
}

func (me *MurmurHash2A) Begin(seed uint32) {
	me.hash, me.tail, me.count, me.size = seed, 0, 0, 0
}

func (me *MurmurHash2A) Add(data []byte) {
	i := 0
	length := len(data)
	me.size += uint32(length)

	me.mixTail(data, &i, &length)

	var k uint32
	for length >= 4 {
		k = binary.LittleEndian.Uint32(data[i:])
		me.hash, k = mmix(me.hash, k)
		i += 4
		length -= 4
	}

	me.mixTail(data, &i, &length)
}

func (me *MurmurHash2A) End() uint32 {
	me.hash, me.tail = mmix(me.hash, me.tail)
	me.hash, me.size = mmix(me.hash, me.size)

	me.hash ^= me.hash >> 13
	me.hash *= m
	me.hash ^= me.hash >> 15

	return me.hash
}

func New() hash.Hash32 {
	me := new(MurmurHash2A)
	me.Reset()
	return me
}
func (me *MurmurHash2A) Reset() { me.Begin(0) }
func (me *MurmurHash2A) Size() int {
	return Size
}
func (me *MurmurHash2A) Write(p []byte) (n int, err os.Error) {
	me.Add(p)
	return len(p), nil
}
func (me *MurmurHash2A) Sum32() uint32 {
	return me.End()
}

func (me *MurmurHash2A) Sum() []byte {
	p := make([]byte, 4)
	s := me.Sum32()
	binary.BigEndian.PutUint32(p, s)
	return p
}

func main() {
	flag.Parse()
	me := new(MurmurHash2A)
	me.Begin(0)
	for i := 0; i < flag.NArg(); i++ {
		//me.Begin(0)
		s := flag.Arg(i)
		me.Add([]byte(s))
		//h := me.End()
		//fmt.Printf("%s: %d (%x)\n", s, h, h)
	}
	h := me.End()
	fmt.Printf("%d (%x)\n", h, h)
}
