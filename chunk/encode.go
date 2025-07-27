package chunk

import (
	"bytes"
)

const (
	// SubChunkVersion is the current version of the written sub chunks, specifying the format they are
	// written on disk and over network.
	SubChunkVersion = 9
	// CurrentBlockVersion is the current version of blocks (states) of the game. This version is composed
	// of 4 bytes indicating a version, interpreted as a big endian int. The current version represents
	// 1.16.0.14 {1, 16, 0, 14}.
	CurrentBlockVersion int32 = 18040335
)

// EncodeSubChunk encodes a sub-chunk from a chunk into bytes. An Encoding may be passed to encode either for network or
// disk purposed, the most notable difference being that the network encoding generally uses varints and no NBT.
func EncodeSubChunk(c *Chunk, e Encoding, ind int) []byte {
	buf := bytes.NewBuffer(nil)

	s := c.sub[ind]
	_, _ = buf.Write([]byte{SubChunkVersion, byte(len(s.storages)), uint8(ind + (c.r[0] >> 4))})
	for _, storage := range s.storages {
		encodePalettedStorage(buf, storage, nil, e)
	}
	sub := make([]byte, buf.Len())
	_, _ = buf.Read(sub)
	return sub
}

// encodePalettedStorage encodes a PalettedStorage into a bytes.Buffer. The Encoding passed is used to write the Palette
// of the PalettedStorage.
func encodePalettedStorage(buf *bytes.Buffer, storage, previous *PalettedStorage, e Encoding) {
	if storage.Equal(previous) {
		_, _ = buf.Write([]byte{0x7f<<1 | e.network()})
		return
	}
	b := make([]byte, len(storage.indices)*4+1)
	b[0] = byte(storage.bitsPerIndex<<1) | e.network()

	for i, v := range storage.indices {
		// Explicitly don't use the binary package to greatly improve performance of writing the uint32s.
		b[i*4+1], b[i*4+2], b[i*4+3], b[i*4+4] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24)
	}
	_, _ = buf.Write(b)

	e.encodePalette(buf, storage.palette)
}
