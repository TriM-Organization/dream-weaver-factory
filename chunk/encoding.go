package chunk

import (
	"bytes"
	"fmt"

	block_general "github.com/TriM-Organization/dream-weaver-factory/block/general"
	"github.com/TriM-Organization/dream-weaver-factory/convert"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	UpgradeTypeNone = iota
	UpgradeTypeToNewer
	UpgradeTypeToLower
)

type (
	// Encoding is an encoding type used for Chunk encoding. Implementations of this interface are DiskEncoding and
	// NetworkEncoding, which can be used to encode a Chunk to an intermediate disk or network representation respectively.
	Encoding interface {
		encodePalette(buf *bytes.Buffer, p *Palette)
		decodePalette(buf *bytes.Buffer, blockSize paletteSize, upgradeType int) (*Palette, error)
		network() byte
	}
)

var (
	// NetworkEncoding is the Encoding used for sending a Chunk over network. It does not use NBT and writes varints.
	NetworkEncoding networkEncoding
)

// networkEncoding implements the Chunk encoding for sending over network.
type networkEncoding struct{}

func (networkEncoding) network() byte { return 1 }

func (networkEncoding) encodePalette(buf *bytes.Buffer, p *Palette) {
	if p.size != 0 {
		_ = protocol.WriteVarint32(buf, int32(p.Len()))
	}
	for _, val := range p.values {
		_ = protocol.WriteVarint32(buf, int32(val))
	}
}

func (networkEncoding) decodePalette(buf *bytes.Buffer, blockSize paletteSize, upgradeType int) (*Palette, error) {
	var paletteCount int32 = 1
	if blockSize != 0 {
		if err := protocol.Varint32(buf, &paletteCount); err != nil {
			return nil, fmt.Errorf("error reading palette entry count: %w", err)
		}
		if paletteCount <= 0 {
			return nil, fmt.Errorf("invalid palette entry count %v", paletteCount)
		}
	}

	blocks, temp := make([]uint32, paletteCount), int32(0)
	for i := int32(0); i < paletteCount; i++ {
		var rid uint32
		var found bool

		if err := protocol.Varint32(buf, &temp); err != nil {
			return nil, fmt.Errorf("error decoding palette entry: %w", err)
		}

		switch upgradeType {
		case UpgradeTypeToNewer:
			rid, found = convert.NEMCToStdBlockRuntimeID(uint32(temp))
			if !found {
				rid, _ = block_general.StdStateToRuntimeID("minecraft:unknown", map[string]any{})
			}
		case UpgradeTypeToLower:
			rid, found = convert.StdToNEMCBlockRuntimeID(uint32(temp))
			if !found {
				rid, _ = block_general.NEMCStateToRuntimeID("minecraft:unknown", map[string]any{})
			}
		default:
			rid = uint32(temp)
		}

		blocks[i] = rid
	}

	return &Palette{values: blocks, size: blockSize}, nil
}
