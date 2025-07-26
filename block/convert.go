package block

// NEMCToStdBlockRuntimeID ..
func NEMCToStdBlockRuntimeID(rid uint32) (result uint32, found bool) {
	blockEntry, ok := olderToNewerBlock[rid]
	if ok {
		return blockEntry.RuntimeID, true
	}
	return 0, false
}

// StdToNEMCBlockRuntimeID ..
func StdToNEMCBlockRuntimeID(rid uint32) (result uint32, found bool) {
	blockEntry, ok := newerToOlderBlock[rid]
	if ok {
		return blockEntry.RuntimeID, true
	}
	return 0, false
}
