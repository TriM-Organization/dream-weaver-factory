package convert

import (
	block_general "github.com/TriM-Organization/dream-weaver-factory/block/general"
	_ "github.com/TriM-Organization/dream-weaver-factory/block/nemc"
	_ "github.com/TriM-Organization/dream-weaver-factory/block/std"
	"github.com/TriM-Organization/dream-weaver-factory/define"
	"github.com/TriM-Organization/dream-weaver-factory/upgrader/blockupgrader"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

var (
	olderToNewerBlock map[uint32]define.BlockEntry
	newerToOlderBlock map[uint32]define.BlockEntry
)

func init() {
	olderToNewerBlock = make(map[uint32]define.BlockEntry)
	newerToOlderBlock = make(map[uint32]define.BlockEntry)

	for _, blockEntry := range block_general.NEMCAllBlocks {
		// Upgrade current block
		newerBlockState := blockupgrader.Upgrade(blockupgrader.BlockState{
			Name:       blockEntry.Block.Name,
			Properties: deepCopyBlockStates(blockEntry.Block.Properties),
			Version:    blockEntry.Block.Version,
		})
		newerRuntimeID := block_general.ComputeBlockHash(
			newerBlockState.Name,
			newerBlockState.Properties,
		)

		// Set mapping
		olderToNewerBlock[blockEntry.RuntimeID] = define.BlockEntry{
			Block: define.BlockState{
				Name:       newerBlockState.Name,
				Properties: newerBlockState.Properties,
				Version:    newerBlockState.Version,
			},
			RuntimeID: newerRuntimeID,
		}
		if _, ok := newerToOlderBlock[newerRuntimeID]; !ok {
			newerToOlderBlock[newerRuntimeID] = blockEntry
		}

		// Special fix for some specific blocks
		if blockEntry.Block.Name == "minecraft:brown_mushroom_block" {
			hugeMushroomBits, _ := blockEntry.Block.Properties["huge_mushroom_bits"].(int32)
			if hugeMushroomBits == 10 || hugeMushroomBits == 15 {
				_, ok := newerToOlderBlock[newerRuntimeID]
				if !ok {
					panic("init: Should nerver happened")
				}
				newerToOlderBlock[newerRuntimeID] = blockEntry
			}
		}
	}
}

// deepCopyBlockStates ..
func deepCopyBlockStates(src map[string]any) (dst map[string]any) {
	nbtBytes, _ := nbt.MarshalEncoding(src, nbt.LittleEndian)
	_ = nbt.UnmarshalEncoding(nbtBytes, &dst, nbt.LittleEndian)
	return
}
