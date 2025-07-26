package block_general

import "github.com/TriM-Organization/dream-weaver-factory/define"

const UseNetworkBlockRuntimeID = true

var (
	// NEMCAirRuntimeID is the runtime ID of an air block in netease Minecraft.
	NEMCAirRuntimeID uint32
	// StdAirRuntimeID is the runtime ID of an air block in standard Minecraft.
	StdAirRuntimeID uint32
)

var (
	// NEMCAllBlocks holds all blocks for netease Minecraft.
	NEMCAllBlocks []define.BlockEntry
	// StdAllBlocks holds all blocks for standard Minecraft.
	StdAllBlocks []define.BlockEntry
)

var (
	// NEMCRuntimeIDToState must hold a function to convert a runtime
	// ID to a name and its state properties that align with netease Minecraft.
	NEMCRuntimeIDToState func(runtimeID uint32) (name string, properties map[string]any, found bool)
	// NEMCStateToRuntimeID must hold a function to convert a name and
	// its state properties to a runtime ID that align with netease Minecraft.
	NEMCStateToRuntimeID func(name string, properties map[string]any) (runtimeID uint32, found bool)
)

var (
	// StdRuntimeIDToState must hold a function to convert a runtime
	// ID to a name and its state properties that align with standard Minecraft.
	StdRuntimeIDToState func(runtimeID uint32) (name string, properties map[string]any, found bool)
	// StdStateToRuntimeID must hold a function to convert a name and
	// its state properties to a runtime ID that align with standard Minecraft.
	StdStateToRuntimeID func(name string, properties map[string]any) (runtimeID uint32, found bool)
)
