package chunk

import (
	"bytes"
	"fmt"

	block_general "github.com/TriM-Organization/dream-weaver-factory/block/general"
	"github.com/TriM-Organization/dream-weaver-factory/define"
)

// DecodeSubChunk decodes a SubChunk from a bytes.Buffer.
// The Encoding passed defines how the block storages of the
// SubChunk are decoded.
func DecodeSubChunk(buf *bytes.Buffer, r define.Range, e Encoding, upgradeType int) (subChunk *SubChunk, index int, err error) {
	var airRuntimeID uint32
	switch upgradeType {
	case UpgradeTypeToLower:
		airRuntimeID = block_general.NEMCAirRuntimeID
	default:
		airRuntimeID = block_general.StdAirRuntimeID
	}
	ver, err := buf.ReadByte()
	if err != nil {
		return nil, 255, fmt.Errorf("error reading version: %w", err)
	}
	sub := NewSubChunk(airRuntimeID)
	switch ver {
	default:
		return nil, 255, fmt.Errorf("unknown sub chunk version %v: can't decode", ver)
	case 1:
		// Version 1 only has one layer for each sub chunk, but uses the format with palettes.
		storage, err := decodePalettedStorage(buf, e, upgradeType)
		if err != nil {
			return nil, 255, err
		}
		sub.storages = append(sub.storages, storage)
	case 8, 9:
		// Version 8 allows up to 256 layers for one sub chunk.
		storageCount, err := buf.ReadByte()
		if err != nil {
			return nil, 255, fmt.Errorf("error reading storage count: %w", err)
		}
		if ver == 9 {
			uIndex, err := buf.ReadByte()
			if err != nil {
				return nil, 255, fmt.Errorf("error reading sub-chunk index: %w", err)
			}
			// The index as written here isn't the actual index of the sub-chunk within the chunk. Rather, it is the Y
			// value of the sub-chunk. This means that we need to translate it to an index.
			index = int(int8(uIndex) - int8(r[0]>>4))
		}
		sub.storages = make([]*PalettedStorage, storageCount)

		for i := range storageCount {
			sub.storages[i], err = decodePalettedStorage(buf, e, upgradeType)
			if err != nil {
				return nil, 255, err
			}
		}
	}
	return sub, index, nil
}

// decodePalettedStorage decodes a PalettedStorage from a bytes.Buffer. The Encoding passed is used to read either a
// network or disk block storage.
func decodePalettedStorage(buf *bytes.Buffer, e Encoding, upgradeType int) (*PalettedStorage, error) {
	blockSize, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading block size: %w", err)
	}
	blockSize >>= 1
	if blockSize == 0x7f {
		return nil, nil
	}

	size := paletteSize(blockSize)
	if size > 32 {
		return nil, fmt.Errorf("cannot read paletted storage (size=%v): size too large", blockSize)
	}
	uint32Count := size.uint32s()

	uint32s := make([]uint32, uint32Count)
	byteCount := uint32Count * 4

	data := buf.Next(byteCount)
	if len(data) != byteCount {
		return nil, fmt.Errorf("cannot read paletted storage (size=%v): not enough block data present: expected %v bytes, got %v", blockSize, byteCount, len(data))
	}
	for i := 0; i < uint32Count; i++ {
		// Explicitly don't use the binary package to greatly improve performance of reading the uint32s.
		uint32s[i] = uint32(data[i*4]) | uint32(data[i*4+1])<<8 | uint32(data[i*4+2])<<16 | uint32(data[i*4+3])<<24
	}
	p, err := e.decodePalette(buf, paletteSize(blockSize), upgradeType)
	return newPalettedStorage(uint32s, p), err
}
