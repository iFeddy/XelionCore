package database

import (
	"fmt"
	"strings"
)

func (r *MemoryCharacterRepository) RenameGuild(guildID int, newName string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	trimmed := strings.TrimSpace(newName)
	if guildID <= 0 {
		return -1311, nil
	}
	if trimmed == "" {
		return -1312, nil
	}
	g, ok := r.guilds[guildID]
	if !ok || g == nil {
		return -1321, nil
	}
	for id, other := range r.guilds {
		if id != guildID && other != nil && strings.EqualFold(other.Name, trimmed) {
			return -1322, nil
		}
	}
	g.Name = trimmed
	return 0, nil
}

func (r *MemoryCharacterRepository) RetypeGuild(guildID int, newType byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.guilds[guildID]
	if !ok || g == nil {
		return 0, nil
	}
	g.Type = newType
	return 1, nil
}

func (r *MemoryCharacterRepository) SaveMysteryBolt(charNo int, mysteryType byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.characters[charNo]; !ok {
		return -3192, nil
	}
	if r.options[charNo] == nil {
		r.options[charNo] = &CharacterOptions{CharNo: charNo}
	}
	return 0, nil
}

func (r *MemoryCharacterRepository) SaveLinkScrollDestMap(charNo int, destMap string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.characters[charNo]
	if !ok || ch == nil {
		return fmt.Errorf("character %d not found", charNo)
	}
	ch.KQMapName = strings.TrimSpace(destMap)
	return nil
}

func (r *MemoryCharacterRepository) UseCharacterChangeIDItem(charNo int, bInit byte, userNo int, itemKey int64, itemOption int16, itemLot int16) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.characters[charNo]
	if !ok || ch == nil {
		return -3023, nil
	}
	for _, entry := range r.changeIDs[userNo] {
		if entry != nil && entry.CharNo == charNo {
			return -3022, nil
		}
	}
	var found *Item
	for i, it := range r.items[charNo] {
		if it != nil && it.ItemKey == itemKey {
			found = it
			if itemLot <= 0 {
				r.items[charNo] = append(r.items[charNo][:i], r.items[charNo][i+1:]...)
			} else {
				it.Flags = int32(itemLot)
			}
			break
		}
	}
	if found == nil {
		return -3036, nil
	}
	ch.ChangeID = &CharIDChangeInfo{
		RowNo:        uint32(len(r.changeIDs[userNo]) + 1),
		UserNo:       userNo,
		CharNo:       charNo,
		NeedChangeID: 1,
		Init:         bInit,
	}
	r.changeIDs[userNo] = append(r.changeIDs[userNo], ch.ChangeID)
	return 0, nil
}

func (r *MemoryCharacterRepository) SetAdminLevel(charNo int, level byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.characters[charNo]
	if !ok || ch == nil {
		return 0, nil
	}
	ch.AdminLevel = level
	return 1, nil
}

func (r *MemoryCharacterRepository) CreateGMItemRewardGroup(operID, ip, cause string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextItemKey++
	return r.nextItemKey, nil
}

func (r *MemoryCharacterRepository) GrantGMItemReward(groupNo int64, charNo int, storageType byte, itemID int, itemKey int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.characters[charNo]; !ok {
		return -1111, fmt.Errorf("character %d not found", charNo)
	}
	if itemKey == 0 {
		r.nextItemKey++
		itemKey = r.nextItemKey
	}
	for _, it := range r.items[charNo] {
		if it != nil && it.StorageType == storageType {
			it.Storage++
		}
	}
	r.items[charNo] = append(r.items[charNo], &Item{
		ItemKey:     itemKey,
		StorageType: storageType,
		Owner:       charNo,
		Storage:     0,
		ItemID:      uint16(itemID),
		Flags:       1,
	})
	return 0, nil
}
