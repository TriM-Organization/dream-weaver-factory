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
		if newerBlockState.Name == "minecraft:micro_block" {
			continue
		}
		if newerBlockState.Name == "minecraft:mod_ore" {
			continue
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

		// Special fix for some specific blocks
		if blockEntry.Block.Name == "minecraft:brown_mushroom_block" {
			hugeMushroomBits, _ := blockEntry.Block.Properties["huge_mushroom_bits"].(int32)
			if hugeMushroomBits == 10 || hugeMushroomBits == 15 {
				_, ok := newerToOlderBlock[newerRuntimeID]
				if !ok {
					panic("init: Should never happened")
				}
				newerToOlderBlock[newerRuntimeID] = blockEntry
			}
		}
	}

	monkeyFix()
}

// deepCopyBlockStates ..
func deepCopyBlockStates(src map[string]any) (dst map[string]any) {
	nbtBytes, _ := nbt.MarshalEncoding(src, nbt.LittleEndian)
	_ = nbt.UnmarshalEncoding(nbtBytes, &dst, nbt.LittleEndian)
	return
}

// monkeyFix (Fix by hand and eyes)
func monkeyFix() {
	// Skull block
	{
		skullBlockNames := []string{
			"minecraft:skeleton_skull",
			"minecraft:wither_skeleton_skull",
			"minecraft:zombie_head",
			"minecraft:player_head",
			"minecraft:creeper_head",
			"minecraft:dragon_head",
			"minecraft:piglin_head",
		}

		for _, blockname := range skullBlockNames {
			for facingDirection := range int32(6) {
				blockState := map[string]any{
					"facing_direction": facingDirection,
				}

				olderRuntimeID, found := block_general.NEMCStateToRuntimeID("minecraft:skull", blockState)
				if !found {
					panic("monkeyFix: Should never happened")
				}
				newerRuntimeID, found := block_general.StdStateToRuntimeID(blockname, blockState)
				if !found {
					panic("monkeyFix: Should never happened")
				}

				newerToOlderBlock[newerRuntimeID] = define.BlockEntry{
					Block: define.BlockState{
						Name:       "minecraft:skull",
						Properties: blockState,
					},
					RuntimeID: olderRuntimeID,
				}
			}
		}
	}
}
