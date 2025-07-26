package convert

import "github.com/df-mc/worldupgrader/itemupgrader"

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

// UpgradeItem ..
func UpgradeItem(name string, meta int16) (newName string, newMeta int16) {
	newItem := itemupgrader.Upgrade(itemupgrader.ItemMeta{
		Name: name,
		Meta: meta,
	})
	return newItem.Name, newItem.Meta
}
