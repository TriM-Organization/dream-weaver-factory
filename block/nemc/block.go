package block_nemc

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"
	"io"
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
	type nemc struct {
		Blocks []define.NetEaseBlock `nbt:"blocks"`
	}

	var neteaseBlocks nemc
	gzipReader, err := gzip.NewReader(bytes.NewBuffer(blockStates))
	if err != nil {
		panic(`init: Failed to unzip "block_states.nbt" (Stage 1)`)
	}
	defer gzipReader.Close()

	unzipedBytes, err := io.ReadAll(gzipReader)
	if err != nil {
		panic(`init: Failed to unzip "block_states.nbt" (Stage 2)`)
	}

	err = nbt.NewDecoderWithEncoding(bytes.NewBuffer(unzipedBytes), nbt.BigEndian).Decode(&neteaseBlocks)
	if err != nil {
		panic("init: Failed to decode netease blocks from NBT")
	}

	// Register all block states present in the block_states.nbt file. These are all possible options registered
	// blocks may encode to.
	for _, value := range neteaseBlocks.Blocks {
		registerBlockState(define.BlockState{
			Name:       value.Name,
			Properties: value.States,
			Version:    value.Version,
		})
	}

	block_general.NEMCRuntimeIDToState = func(runtimeID uint32) (name string, properties map[string]any, found bool) {
		blockEntry, found := blockStateMapping[runtimeID]
		if found {
			return blockEntry.Block.Name, blockEntry.Block.Properties, true
		}
		return "", nil, false
	}
	block_general.NEMCStateToRuntimeID = func(name string, properties map[string]any) (runtimeID uint32, found bool) {
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

	if block_general.NEMCUseBlockNetworkIDHashes {
		rid = hash
	} else {
		rid = uint32(len(blockStateMapping))
	}
	if s.Name == "minecraft:air" {
		block_general.NEMCAirRuntimeID = rid
	}

	blockEntry := define.BlockEntry{
		Block:     s,
		RuntimeID: rid,
	}
	blockStateMapping[hash] = blockEntry
	block_general.NEMCAllBlocks = append(block_general.NEMCAllBlocks, blockEntry)
}
