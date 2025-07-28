package block_std

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	block_general "github.com/TriM-Organization/dream-weaver-factory/block/general"
	"github.com/TriM-Organization/dream-weaver-factory/define"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

var (
	//go:embed block_states.nbt
	blockStates []byte

	// blockProperties ..
	blockProperties = map[string]map[string]any{}
	// blockStateMapping holds a map for looking up a block entry by the network runtime id it produces.
	blockStateMapping = map[uint32]define.BlockEntry{}
)

func init() {
	// Register all block states present in the block_states.nbt file. These are all possible options registered
	// blocks may encode to.
	dec := nbt.NewDecoder(bytes.NewBuffer(blockStates))
	for {
		var s define.BlockState
		if err := dec.Decode(&s); err != nil {
			break
		}
		registerBlockState(s)
	}

	block_general.StdRuntimeIDToState = func(runtimeID uint32) (name string, properties map[string]any, found bool) {
		blockEntry, found := blockStateMapping[runtimeID]
		if found {
			return blockEntry.Block.Name, blockEntry.Block.Properties, true
		}
		return "", nil, false
	}
	block_general.StdStateToRuntimeID = func(name string, properties map[string]any) (runtimeID uint32, found bool) {
		if !strings.HasPrefix(name, "minecraft:") {
			name = "minecraft:" + name
		}

		networkRuntimeID := block_general.ComputeBlockHash(name, properties)
		if blockEntry, ok := blockStateMapping[networkRuntimeID]; ok {
			return blockEntry.RuntimeID, true
		}

		networkRuntimeID = block_general.ComputeBlockHash(name, blockProperties[name])
		blockEntry, ok := blockStateMapping[networkRuntimeID]
		return blockEntry.RuntimeID, ok
	}
}

// registerBlockState registers a new blockState to the states slice.
// The function panics if the properties the blockState hold are invalid
// or if the blockState was already registered.
func registerBlockState(s define.BlockState) {
	var rid uint32

	hash := block_general.ComputeBlockHash(s.Name, s.Properties)
	if _, ok := blockStateMapping[hash]; ok {
		panic(fmt.Sprintf("cannot register the same state twice (%+v)", s))
	}
	if _, ok := blockProperties[s.Name]; !ok {
		blockProperties[s.Name] = s.Properties
	}

	if block_general.StdUseBlockNetworkIDHashes {
		rid = hash
	} else {
		rid = uint32(len(blockStateMapping))
	}
	if s.Name == "minecraft:air" {
		block_general.StdAirRuntimeID = rid
	}

	blockEntry := define.BlockEntry{
		Block:     s,
		RuntimeID: rid,
	}
	blockStateMapping[hash] = blockEntry
	block_general.StdAllBlocks = append(block_general.StdAllBlocks, blockEntry)
}
