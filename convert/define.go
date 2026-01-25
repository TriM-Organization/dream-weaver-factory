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
		var shouldUpgrade bool
		var newerBlockState blockupgrader.BlockState

		// Special fix for some specific blocks
		switch blockEntry.Block.Name {
		case "minecraft:red_mushroom_block", "minecraft:brown_mushroom_block", "minecraft:mushroom_stem":
			hugeMushroomBits, _ := blockEntry.Block.Properties["huge_mushroom_bits"].(int32)
			if hugeMushroomBits != 10 && hugeMushroomBits != 15 {
				shouldUpgrade = true
			}
		default:
			shouldUpgrade = true
		}

		// Upgrade current block
		if shouldUpgrade {
			newerBlockState = blockupgrader.Upgrade(blockupgrader.BlockState{
				Name:       blockEntry.Block.Name,
				Properties: deepCopyBlockStates(blockEntry.Block.Properties),
				Version:    blockEntry.Block.Version,
			})
			if newerBlockState.Name == "minecraft:micro_block" {
				continue
			}
		} else {
			newerBlockState = blockupgrader.BlockState{
				Name:       blockEntry.Block.Name,
				Properties: deepCopyBlockStates(blockEntry.Block.Properties),
				Version:    blockEntry.Block.Version,
			}
		}

		newerRuntimeID, found := block_general.StdStateToRuntimeID(
			newerBlockState.Name,
			newerBlockState.Properties,
		)
		if !found {
			panic("init: Should never happened")
		}

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
	}
}

// deepCopyBlockStates ..
func deepCopyBlockStates(src map[string]any) (dst map[string]any) {
	nbtBytes, _ := nbt.MarshalEncoding(src, nbt.LittleEndian)
	_ = nbt.UnmarshalEncoding(nbtBytes, &dst, nbt.LittleEndian)
	return
}
